package anime

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/mappers"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrAnimeNotFound     = errors.New("anime not found")
	ErrCharacterNotFound = errors.New("character not found")
	ErrPersonNotFound    = errors.New("person not found")
	TrailerNotFound      = errors.New("trailer not found")
	BannerNotFound       = errors.New("banner not found")
)

func (s *AnimeService) GetAnimeByID(
	ctx context.Context,
	id string,
) (*models.AnimeWithMetadataResponse, error) {
	a, err := s.repo.GetAnimeById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAnimeNotFound
	}
	if err != nil {
		return nil, err
	}

	m, err := s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
	switch {
	case err == nil && time.Since(m.UpdatedAt.Time) < s.refresher.ttl:
		dto := mappers.AnimeWithMetadataFromRepository(a, m)
		return &dto, nil

	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	// if error and no existing metadata return error
	if err := s.refresher.RefreshBlocking(ctx, a.MalID.Int32); err != nil && m.MalID == 0 {
		return nil, err
	}

	m, err = s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
	if err != nil {
		return nil, err
	}

	dto := mappers.AnimeWithMetadataFromRepository(a, m)

	go func() {
		if a.Season == m.Season && a.SeasonYear == m.SeasonYear.Int32 {
			return
		}

		s.repo.UpdateAnimeSeasons(context.Background(), repository.UpdateAnimeSeasonsParams{
			Season:     m.Season,
			SeasonYear: m.SeasonYear.Int32,
		})
	}()

	return &dto, nil
}

func (s *AnimeService) GetAnimeTrailer(ctx context.Context, id string) (*models.TrailerResponse, error) {
	a, err := s.GetAnimeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if a.Metadata.TrailerEmbedURL != "" {
		return &models.TrailerResponse{Trailer: a.Metadata.TrailerEmbedURL}, nil
	}

	malID := *a.MalID
	if malID == 0 {
		return nil, ErrAnimeNotFound
	}

	t, err := s.malClient.GetTrailer(ctx, int(malID))
	if err != nil || t == "" {
		return nil, TrailerNotFound
	}

	a.Metadata.TrailerEmbedURL = t
	params := repository.UpdateAnimeMetadataTrailerParams{
		TrailerEmbedUrl: pgtype.Text{String: t, Valid: true},
		MalID:           malID,
	}
	if err := s.repo.UpdateAnimeMetadataTrailer(ctx, params); err != nil {
		return nil, fmt.Errorf("failed to update metadata for MAL ID %d: %v", malID, err)
	}

	return &models.TrailerResponse{Trailer: t}, nil
}

func (s *AnimeService) GetAnimeBanner(ctx context.Context, id string) (models.BannerResponse, error) {
	banner, err := cache.GetOrFill(ctx, s.redis, fmt.Sprintf("anime_banner:%s", id), 30*24*time.Hour, func(ctx context.Context) (string, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrAnimeNotFound
		}
		if err != nil {
			return "", fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		anime, err := s.anilistClient.GetAnimeDetails(ctx, int(a.MalID.Int32))
		if err != nil {
			return "", fmt.Errorf("failed to fetch anime details from Anilist for MAL ID %d: %v", a.MalID.Int32, err)
		}

		if anime.Media.BannerImage == "" {
			return "", BannerNotFound
		}

		return anime.Media.BannerImage, nil
	})

	if err != nil {
		return models.BannerResponse{}, err
	}

	return models.BannerResponse{URL: banner}, nil
}

func (s *AnimeService) GetAnimeRelations(ctx context.Context, id string) (models.RelationsResponse, error) {
	key := fmt.Sprintf("anime_relations:%s", id)

	return cache.GetOrFill(ctx, s.redis, key, 7*24*time.Hour, func(ctx context.Context) (models.RelationsResponse, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.RelationsResponse{}, ErrAnimeNotFound
		}
		if err != nil {
			return models.RelationsResponse{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		if !a.MalID.Valid || a.MalID.Int32 <= 0 {
			return models.RelationsResponse{}, fmt.Errorf("invalid MAL ID for anime ID %s", id)
		}

		malID := int(a.MalID.Int32)
		fr, err := s.shikimoriClient.GetAnimeFranchise(ctx, malID)
		if err != nil {
			return models.RelationsResponse{}, fmt.Errorf("failed to fetch franchise for MAL ID %d: %v", malID, err)
		}

		watchIDs := deriveWatchOrder(fr, malID)
		relatedIDs := deriveRelated(fr, malID, watchIDs)
		fullIDs := deriveFullFranchise(fr)

		rows, err := s.repo.GetAnimesByMalIds(ctx, func(ids []int) []int32 {
			out := make([]int32, 0, len(ids))
			for _, id := range ids {
				out = append(out, int32(id))
			}
			return out
		}(fullIDs))
		if err != nil {
			return models.RelationsResponse{}, fmt.Errorf("failed to fetch related animes: %w", err)
		}

		dtoMap := make(map[int32]models.AnimeResponse, len(rows))
		for _, r := range rows {
			dtoMap[r.MalID.Int32] = mappers.AnimeFromRepository(r)
			s.refresher.Enqueue(r.MalID.Int32)
		}

		sliceDto := func(ids []int) []models.AnimeResponse {
			out := make([]models.AnimeResponse, 0, len(ids))
			for _, id := range ids {
				if dto, ok := dtoMap[int32(id)]; ok {
					out = append(out, dto)
				}
			}
			return out
		}

		reverse := func(ids []int) []int {
			out := make([]int, len(ids))
			for i, id := range ids {
				out[len(ids)-1-i] = id
			}
			return out
		}

		var relations models.RelationsResponse

		if len(watchIDs) > 1 && slices.Contains(watchIDs, malID) {
			relations = models.RelationsResponse{
				WatchOrder: sliceDto(watchIDs),
				Related:    sliceDto(reverse(relatedIDs)),
			}
		} else {
			relations = models.RelationsResponse{
				WatchOrder: []models.AnimeResponse{},
				Related: sliceDto(func(ids []int) []int {
					if len(ids) == 1 {
						return []int{}
					}
					return reverse(ids)
				}(fullIDs)),
			}
		}

		return relations, nil
	})
}

func (s *AnimeService) GetAnimeCharacters(ctx context.Context, id string) (models.CharactersResponse, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("anime_characters:%s", id), 7*24*time.Hour, func(ctx context.Context) (models.CharactersResponse, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.CharactersResponse{}, ErrAnimeNotFound
		}
		if err != nil {
			return models.CharactersResponse{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		if !a.MalID.Valid || a.MalID.Int32 <= 0 {
			return models.CharactersResponse{}, fmt.Errorf("invalid MAL ID for anime ID %s", id)
		}

		characters, err := s.jikanClient.GetAnimeCharacters(ctx, int(a.MalID.Int32))
		if err != nil {
			return models.CharactersResponse{}, fmt.Errorf("failed to fetch characters from MAL for MAL ID %d: %v", a.MalID.Int32, err)
		}

		var dto models.CharactersResponse
		for _, c := range characters {
			dto = append(dto, models.CharacterResponse{
				MalID:     int32(c.MalID),
				Name:      c.Name,
				Role:      c.Role,
				Favorites: int32(c.Favorites),
				Image:     c.Image.ImageURL,
			})
		}

		sort.Slice(dto, func(i, j int) bool {
			if dto[i].Role == dto[j].Role {
				return dto[i].Favorites > dto[j].Favorites
			}
			return dto[i].Role < dto[j].Role
		})

		return dto, nil
	})
}

func (s *AnimeService) GetCharacterFull(ctx context.Context, malID int32) (models.CharacterFullResponse, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("character_full:%d", malID), 7*24*time.Hour, func(ctx context.Context) (models.CharacterFullResponse, error) {
		character, err := s.jikanClient.GetCharacterFull(ctx, int(malID))
		if err != nil {
			return models.CharacterFullResponse{}, fmt.Errorf("failed to fetch character from Jikan for MAL ID %d: %v", malID, err)
		}

		malIDs := make([]int32, 0, len(character.Anime))
		for _, anime := range character.Anime {
			malIDs = append(malIDs, int32(anime.Anime.MalID))
		}

		animeMap := make(map[int32]models.AnimeResponse)
		if len(malIDs) > 0 {
			animes, err := s.repo.GetAnimesByMalIds(ctx, malIDs)
			if err != nil {
				return models.CharacterFullResponse{}, fmt.Errorf("failed to fetch animes from database: %v", err)
			}
			for _, anime := range animes {
				animeMap[anime.MalID.Int32] = mappers.AnimeFromRepository(anime)
			}
		}

		var animeList []models.CharacterAnimeResponse
		for _, anime := range character.Anime {
			malID := int32(anime.Anime.MalID)
			animeData, exists := animeMap[malID]
			if exists {
				animeList = append(animeList, models.CharacterAnimeResponse{
					Role:  anime.Role,
					Anime: animeData,
				})
			}
		}

		var voices []models.CharacterVoiceResponse
		for _, voice := range character.Voices {
			var image *string
			if voice.Person.Images.JPG.ImageURL != "" {
				image = &voice.Person.Images.JPG.ImageURL
			}
			voices = append(voices, models.CharacterVoiceResponse{
				Person: models.CharacterVoicePersonResponse{
					MalID: int32(voice.Person.MalID),
					Name:  voice.Person.Name,
					Image: image,
				},
				Language: voice.Language,
			})
		}

		imageURL := character.Images.Webp.ImageURL
		if imageURL == "" {
			imageURL = character.Images.JPG.ImageURL
		}

		return models.CharacterFullResponse{
			MalID:     int32(character.MalID),
			Name:      character.Name,
			NameKanji: &character.NameKanji,
			Nicknames: character.Nicknames,
			Favorites: int32(character.Favorites),
			About:     &character.About,
			Image:     imageURL,
			Anime:     animeList,
			Voices:    voices,
		}, nil
	})
}

func (s *AnimeService) GetPersonFull(ctx context.Context, malID int32) (models.PersonFullResponse, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("person:%d", malID), 7*24*time.Hour, func(ctx context.Context) (models.PersonFullResponse, error) {
		person, err := s.jikanClient.GetPersonFull(ctx, int(malID))
		if err != nil {
			return models.PersonFullResponse{}, ErrPersonNotFound
		}

		var animeMalIDs []int32
		malIDSet := make(map[int32]bool)

		for _, anime := range person.Anime {
			malID := int32(anime.Anime.MalID)
			if !malIDSet[malID] {
				animeMalIDs = append(animeMalIDs, malID)
				malIDSet[malID] = true
			}
		}

		for _, voice := range person.Voices {
			malID := int32(voice.Anime.MalID)
			if !malIDSet[malID] {
				animeMalIDs = append(animeMalIDs, malID)
				malIDSet[malID] = true
			}
		}

		var animeList []models.PersonAnimeResponse
		var characterList []models.PersonCharacterResponse
		animeMap := make(map[int32]models.AnimeResponse)

		if len(animeMalIDs) > 0 {
			animes, err := s.repo.GetAnimesByMalIds(ctx, animeMalIDs)
			if err != nil {
				return models.PersonFullResponse{}, err
			}

			for _, anime := range animes {
				animeMap[anime.MalID.Int32] = mappers.AnimeFromRepository(anime)
			}

			for _, anime := range person.Anime {
				if animeResp, exists := animeMap[int32(anime.Anime.MalID)]; exists {
					animeList = append(animeList, models.PersonAnimeResponse{
						Position: anime.Position,
						Anime:    animeResp,
					})
				}
			}

			for _, voice := range person.Voices {
				if animeResp, exists := animeMap[int32(voice.Anime.MalID)]; exists {
					characterList = append(characterList, models.PersonCharacterResponse{
						Role:  voice.Role,
						Anime: animeResp,
						Character: models.CharacterResponse{
							MalID:     int32(voice.Character.MalID),
							Name:      voice.Character.Name,
							Role:      voice.Role,
							Favorites: 0,
							Image:     voice.Character.Images.Webp.ImageURL,
						},
					})
				}
			}
		}

		imageURL := person.Images.JPG.ImageURL
		if imageURL == "" {
			imageURL = "https://via.placeholder.com/300x400?text=No+Image"
		}

		return models.PersonFullResponse{
			MalID:          int32(person.MalID),
			Name:           person.Name,
			GivenName:      person.GivenName,
			FamilyName:     person.FamilyName,
			AlternateNames: person.AlternateNames,
			Birthday:       person.Birthday,
			Favorites:      int32(person.Favorites),
			About:          person.About,
			Image:          imageURL,
			Anime:          animeList,
			Characters:     characterList,
		}, nil
	})
}
