package models

type GenresMode string

const (
	AnyGenresMode GenresMode = "any"
	AllGenresMode GenresMode = "all"
)

type SortBy string

const (
	SortByEname     SortBy = "ename"
	SortByJname     SortBy = "jname"
	SortBySeason    SortBy = "season"
	SortByYear      SortBy = "year"
	SortByRelevance SortBy = "relevance"
	SortByUpdatedAt SortBy = "updated_at"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type GetAnimesParams struct {
	Page       int
	Size       int
	Search     *string
	Genres     []string
	GenresMode GenresMode
	Seasons    []string
	Years      []int
	YearMin    *int
	YearMax    *int
	SortBy     SortBy
	SortOrder  SortOrder
}
