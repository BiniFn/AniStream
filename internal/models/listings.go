package models

import (
	"fmt"
	"strings"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/ggicci/httpin/core"
	"github.com/jackc/pgx/v5/pgtype"
)

type GenresMode string

const (
	AnyGenresMode GenresMode = "any"
	AllGenresMode GenresMode = "all"
)

func (m GenresMode) IsValid() bool {
	switch m {
	case AnyGenresMode, AllGenresMode:
		return true
	default:
		return false
	}
}

func (m GenresMode) ToString() (string, error) {
	return string(m), nil
}

func (m *GenresMode) FromString(s string) error {
	switch s {
	case "any":
		*m = AnyGenresMode
	case "all":
		*m = AllGenresMode
	default:
		return fmt.Errorf("invalid GenresMode: %s", s)
	}
	return nil
}

type SortBy string

const (
	SortByEname            SortBy = "ename"
	SortByJname            SortBy = "jname"
	SortBySeason           SortBy = "season"
	SortByYear             SortBy = "year"
	SortByRelevance        SortBy = "relevance"
	SortByUpdatedAt        SortBy = "updated_at"
	SortByAnimeUpdatedAt   SortBy = "anime_updated_at"
	SortByLibraryUpdatedAt SortBy = "library_updated_at"
)

func (s SortBy) IsValid() bool {
	switch s {
	case SortByEname, SortByJname, SortBySeason, SortByYear, SortByRelevance, SortByUpdatedAt, SortByAnimeUpdatedAt, SortByLibraryUpdatedAt:
		return true
	default:
		return false
	}
}

func (s SortBy) ToString() (string, error) {
	return string(s), nil
}

func (s *SortBy) FromString(str string) error {
	switch str {
	case "ename":
		*s = SortByEname
	case "jname":
		*s = SortByJname
	case "season":
		*s = SortBySeason
	case "year":
		*s = SortByYear
	case "relevance":
		*s = SortByRelevance
	case "updated_at":
		*s = SortByUpdatedAt
	case "anime_updated_at":
		*s = SortByAnimeUpdatedAt
	case "library_updated_at":
		*s = SortByLibraryUpdatedAt
	default:
		return fmt.Errorf("invalid SortBy: %s", str)
	}
	return nil
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (o SortOrder) IsValid() bool {
	switch o {
	case SortOrderAsc, SortOrderDesc:
		return true
	default:
		return false
	}
}

func (o SortOrder) ToString() (string, error) {
	return string(o), nil
}

func (o *SortOrder) FromString(s string) error {
	switch s {
	case "asc":
		*o = SortOrderAsc
	case "desc":
		*o = SortOrderDesc
	default:
		return fmt.Errorf("invalid SortOrder: %s", s)
	}
	return nil
}

type GetAnimeCatalogParams struct {
	Page          int        `in:"query=page;default=1"`
	ItemsPerPage  int        `in:"query=itemsPerPage;default=30"`
	Search        *string    `in:"query=search"`
	Genres        []string   `in:"query=genres"`
	GenresMode    GenresMode `in:"query=genresMode"`
	Seasons       []string   `in:"query=seasons"`
	Years         []int      `in:"query=years"`
	YearMin       *int       `in:"query=yearMin"`
	YearMax       *int       `in:"query=yearMax"`
	SortBy        SortBy     `in:"query=sortBy"`
	SortOrder     SortOrder  `in:"query=sortOrder"`
	InLibraryOnly *bool      `in:"query=inLibraryOnly"`
	Status        *string    `in:"query=status"`
}

func (p GetAnimeCatalogParams) Normalize() GetAnimeCatalogParams {
	out := p

	if out.Page <= 0 {
		out.Page = 1
	}
	if out.ItemsPerPage <= 0 || out.ItemsPerPage > 100 {
		out.ItemsPerPage = 30
	}

	if out.Search != nil {
		s := strings.TrimSpace(*out.Search)
		if s == "" {
			out.Search = nil
		} else {
			out.Search = &s
		}
	}

	if len(out.Genres) > 0 {
		g := make([]string, 0, len(out.Genres))
		for _, x := range out.Genres {
			x = strings.TrimSpace(x)
			if x != "" {
				g = append(g, x)
			}
		}
		out.Genres = g
	}

	if !out.GenresMode.IsValid() {
		out.GenresMode = AnyGenresMode
	}
	if !out.SortBy.IsValid() {
		if out.Search != nil {
			out.SortBy = SortByRelevance
		} else {
			out.SortBy = SortByUpdatedAt
		}
	}
	if !out.SortOrder.IsValid() {
		out.SortOrder = SortOrderDesc
	}

	if out.YearMin != nil && out.YearMax != nil && *out.YearMin > *out.YearMax {
		*out.YearMin, *out.YearMax = *out.YearMax, *out.YearMin
	}

	return out
}

func (p GetAnimeCatalogParams) toInt32s(in []int) []int32 {
	if in == nil {
		return nil
	}
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}

func textOpt(ptr *string) pgtype.Text {
	if ptr == nil || strings.TrimSpace(*ptr) == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: strings.TrimSpace(*ptr), Valid: true}
}

func textEnum(s string, valid bool) pgtype.Text {
	return pgtype.Text{String: s, Valid: valid}
}

func int4Opt(ptr *int) pgtype.Int4 {
	if ptr == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*ptr), Valid: true}
}

func libraryStatusOpt(ptr *string) repository.NullLibraryStatus {
	if ptr == nil || strings.TrimSpace(*ptr) == "" {
		return repository.NullLibraryStatus{Valid: false}
	}

	status := strings.TrimSpace(*ptr)
	var libStatus repository.LibraryStatus

	switch status {
	case "planning":
		libStatus = repository.LibraryStatusPlanning
	case "watching":
		libStatus = repository.LibraryStatusWatching
	case "completed":
		libStatus = repository.LibraryStatusCompleted
	case "dropped":
		libStatus = repository.LibraryStatusDropped
	case "paused":
		libStatus = repository.LibraryStatusPaused
	default:
		return repository.NullLibraryStatus{Valid: false}
	}

	return repository.NullLibraryStatus{LibraryStatus: libStatus, Valid: true}
}

func (p GetAnimeCatalogParams) ToRepo(limit, offset int32, userID *string) repository.GetAnimeCatalogParams {
	n := p.Normalize()

	var pgUserID pgtype.Text
	if userID != nil {
		pgUserID = pgtype.Text{String: *userID, Valid: true}
	}

	return repository.GetAnimeCatalogParams{
		Limit:         int32(limit),
		Offset:        int32(offset),
		UserID:        pgUserID,
		Search:        textOpt(n.Search),
		Genres:        n.Genres,
		GenresMode:    textEnum(string(n.GenresMode), n.GenresMode.IsValid()),
		Seasons:       n.Seasons,
		Years:         n.toInt32s(n.Years),
		YearMin:       int4Opt(n.YearMin),
		YearMax:       int4Opt(n.YearMax),
		SortBy:        textEnum(string(n.SortBy), n.SortBy.IsValid()),
		SortOrder:     textEnum(string(n.SortOrder), n.SortOrder.IsValid()),
		LibraryStatus: libraryStatusOpt(n.Status),
	}
}

func (p GetAnimeCatalogParams) ToRepoCount(userID *string) repository.GetAnimeCatalogCountParams {
	n := p.Normalize()

	var pgUserID pgtype.Text
	if userID != nil {
		pgUserID = pgtype.Text{String: *userID, Valid: true}
	}

	return repository.GetAnimeCatalogCountParams{
		UserID:        pgUserID,
		Search:        textOpt(n.Search),
		Genres:        n.Genres,
		GenresMode:    textEnum(string(n.GenresMode), n.GenresMode.IsValid()),
		Seasons:       n.Seasons,
		Years:         n.toInt32s(n.Years),
		YearMin:       int4Opt(n.YearMin),
		YearMax:       int4Opt(n.YearMax),
		LibraryStatus: libraryStatusOpt(n.Status),
	}
}

func init() {
	core.RegisterNamedCoder("GenresMode", func(s *string) (core.Stringable, error) {
		return (*GenresMode)(s), nil
	})

	core.RegisterNamedCoder("SortBy", func(s *string) (core.Stringable, error) {
		return (*SortBy)(s), nil
	})

	core.RegisterNamedCoder("SortOrder", func(s *string) (core.Stringable, error) {
		return (*SortOrder)(s), nil
	})
}

type GenrePreview struct {
	Name     string   `json:"name"`
	Previews []string `json:"previews"`
}
