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
	SortByEname     SortBy = "ename"
	SortByJname     SortBy = "jname"
	SortBySeason    SortBy = "season"
	SortByYear      SortBy = "year"
	SortByRelevance SortBy = "relevance"
	SortByUpdatedAt SortBy = "updated_at"
)

func (s SortBy) IsValid() bool {
	switch s {
	case SortByEname, SortByJname, SortBySeason, SortByYear, SortByRelevance, SortByUpdatedAt:
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
	Page         int        `in:"query=page;default=1"`
	ItemsPerPage int        `in:"query=itemsPerPage;default=30"`
	Search       *string    `in:"query=search"`
	Genres       []string   `in:"query=genres"`
	GenresMode   GenresMode `in:"query=genresMode"`
	Seasons      []string   `in:"query=seasons"`
	Years        []int      `in:"query=years"`
	YearMin      *int       `in:"query=yearMin"`
	YearMax      *int       `in:"query=yearMax"`
	SortBy       SortBy     `in:"query=sortBy"`
	SortOrder    SortOrder  `in:"query=sortOrder"`
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

func (p GetAnimeCatalogParams) ToRepo(limit, offset int32) repository.GetAnimeCatalogParams {
	n := p.Normalize()

	return repository.GetAnimeCatalogParams{
		Limit:      int32(limit),
		Offset:     int32(offset),
		Search:     textOpt(n.Search),
		Genres:     n.Genres,
		GenresMode: textEnum(string(n.GenresMode), n.GenresMode.IsValid()),
		Seasons:    n.Seasons,
		Years:      n.toInt32s(n.Years),
		YearMin:    int4Opt(n.YearMin),
		YearMax:    int4Opt(n.YearMax),
		SortBy:     textEnum(string(n.SortBy), n.SortBy.IsValid()),
		SortOrder:  textEnum(string(n.SortOrder), n.SortOrder.IsValid()),
	}
}

func (p GetAnimeCatalogParams) ToRepoCount() repository.GetAnimeCatalogCountParams {
	n := p.Normalize()

	return repository.GetAnimeCatalogCountParams{
		Search:     textOpt(n.Search),
		Genres:     n.Genres,
		GenresMode: textEnum(string(n.GenresMode), n.GenresMode.IsValid()),
		Seasons:    n.Seasons,
		Years:      n.toInt32s(n.Years),
		YearMin:    int4Opt(n.YearMin),
		YearMax:    int4Opt(n.YearMax),
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
