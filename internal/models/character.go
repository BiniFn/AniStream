package models

type CharacterResponse struct {
	MalID     int32  `json:"malId" validate:"required" example:"1"`
	Name      string `json:"name" validate:"required" example:"Eren Yeager"`
	Role      string `json:"role" validate:"required" example:"Main"`
	Favorites int32  `json:"favorites" validate:"required" example:"10000"`
	Image     string `json:"image" validate:"required" example:"https://example.com/character.jpg"`
}

type CharacterFullResponse struct {
	MalID     int32                    `json:"malId" validate:"required" example:"1"`
	Name      string                   `json:"name" validate:"required" example:"Eren Yeager"`
	NameKanji *string                  `json:"nameKanji" example:"エレン・イェーガー"`
	Nicknames []string                 `json:"nicknames" example:"Eren,Tatakae"`
	Favorites int32                    `json:"favorites" validate:"required" example:"10000"`
	About     *string                  `json:"about" example:"A determined young man who seeks freedom..."`
	Image     string                   `json:"image" validate:"required" example:"https://example.com/character.jpg"`
	Anime     []CharacterAnimeResponse `json:"anime" validate:"required"`
	Voices    []CharacterVoiceResponse `json:"voices" validate:"required"`
}

type CharacterAnimeResponse struct {
	Role  string        `json:"role" validate:"required" example:"Main"`
	Anime AnimeResponse `json:"anime" validate:"required"`
}

type CharacterVoiceResponse struct {
	Person   CharacterVoicePersonResponse `json:"person" validate:"required"`
	Language string                       `json:"language" validate:"required" example:"Japanese"`
}

type CharacterVoicePersonResponse struct {
	MalID int32   `json:"malId" validate:"required" example:"123"`
	Name  string  `json:"name" validate:"required" example:"Yuki Kaji"`
	Image *string `json:"image" example:"https://example.com/voice-actor.jpg"`
}

type CharactersResponse []CharacterResponse

type CharacterListResponse = Pagination[CharacterResponse]

type CharacterFullListResponse = Pagination[CharacterFullResponse]

type PersonFullResponse struct {
	MalID          int32                     `json:"malId" validate:"required" example:"35511"`
	Name           string                    `json:"name" validate:"required" example:"Shion Wakayama"`
	GivenName      *string                   `json:"givenName" example:"詩音"`
	FamilyName     *string                   `json:"familyName" example:"若山"`
	AlternateNames []string                  `json:"alternateNames" example:"Wakayama Shion"`
	Birthday       *string                   `json:"birthday" example:"1998-02-10T00:00:00+00:00"`
	Favorites      int32                     `json:"favorites" validate:"required" example:"1502"`
	About          *string                   `json:"about" example:"Birthplace: Chiba Prefecture, Japan..."`
	Image          string                    `json:"image" validate:"required" example:"https://example.com/person.jpg"`
	Anime          []PersonAnimeResponse     `json:"anime" validate:"required"`
	Characters     []PersonCharacterResponse `json:"characters" validate:"required"`
}

type PersonAnimeResponse struct {
	Position string        `json:"position" validate:"required" example:"Main"`
	Anime    AnimeResponse `json:"anime" validate:"required"`
}

type PersonCharacterResponse struct {
	Role      string            `json:"role" validate:"required" example:"Main"`
	Anime     AnimeResponse     `json:"anime" validate:"required"`
	Character CharacterResponse `json:"character" validate:"required"`
}

type PersonFullListResponse = Pagination[PersonFullResponse]
