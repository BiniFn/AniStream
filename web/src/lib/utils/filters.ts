import type { paths } from '$lib/api/openapi';

export type AnimeListingsQuery = paths['/anime/listings']['get']['parameters']['query'];

export interface FilterState {
	search: string;
	genres: string[];
	genresMode: 'any' | 'all';
	seasons: ('winter' | 'spring' | 'summer' | 'fall' | 'unknown')[];
	years: number[];
	yearMin?: number;
	yearMax?: number;
	sortBy:
		| 'ename'
		| 'jname'
		| 'season'
		| 'year'
		| 'relevance'
		| 'updated_at'
		| 'anime_updated_at'
		| 'library_updated_at';
	sortOrder: 'asc' | 'desc';
	itemsPerPage: number;
	page: number;
	inLibraryOnly?: boolean;
	status?: 'watching' | 'completed' | 'planning' | 'dropped' | 'paused';
}

export const defaultFilters: FilterState = {
	search: '',
	genres: [],
	genresMode: 'any',
	seasons: [],
	years: [],
	sortBy: 'relevance',
	sortOrder: 'desc',
	itemsPerPage: 24,
	page: 1,
};

export function filtersToQuery(filters: FilterState): URLSearchParams {
	const params = new URLSearchParams();
	const apiQuery = filtersToApiQuery(filters);
	const entries = Object.entries(
		apiQuery as Record<string, string | number | boolean | string[] | number[]>,
	);

	for (const [key, value] of entries) {
		if (value !== undefined && value !== null) {
			if (Array.isArray(value)) {
				value.forEach((item) => params.append(key, item.toString()));
			} else {
				params.set(key, value.toString());
			}
		}
	}

	return params;
}

export function filtersToApiQuery(filters: FilterState): AnimeListingsQuery {
	const query: AnimeListingsQuery = {};

	if (filters.search) query.search = filters.search;
	if (filters.genres.length > 0) query.genres = filters.genres;
	if (filters.genresMode !== 'any') query.genresMode = filters.genresMode;
	if (filters.seasons.length > 0) query.seasons = filters.seasons;
	if (filters.years.length > 0) query.years = filters.years;
	if (filters.yearMin !== undefined) query.yearMin = filters.yearMin;
	if (filters.yearMax !== undefined) query.yearMax = filters.yearMax;
	if (filters.sortBy !== 'relevance') query.sortBy = filters.sortBy;
	if (filters.sortOrder !== 'desc') query.sortOrder = filters.sortOrder;
	if (filters.itemsPerPage !== 24) query.itemsPerPage = filters.itemsPerPage;
	if (filters.page !== 1) query.page = filters.page;
	if (filters.inLibraryOnly !== undefined) query.inLibraryOnly = filters.inLibraryOnly;
	if (filters.status !== undefined) query.status = filters.status;

	return query;
}

export function searchParamsToFilters(query: URLSearchParams): FilterState {
	return {
		search: query.get('search') || '',
		genres: query.getAll('genres'),
		genresMode: (query.get('genresMode') as 'any' | 'all') || 'any',
		seasons: query.getAll('seasons') as ('winter' | 'spring' | 'summer' | 'fall' | 'unknown')[],
		years: query.getAll('years').map(Number).filter(Boolean),
		yearMin: query.get('yearMin') ? Number(query.get('yearMin')) : undefined,
		yearMax: query.get('yearMax') ? Number(query.get('yearMax')) : undefined,
		sortBy: (query.get('sortBy') as FilterState['sortBy']) || 'relevance',
		sortOrder: (query.get('sortOrder') as 'asc' | 'desc') || 'desc',
		itemsPerPage: Number(query.get('itemsPerPage')) || 24,
		page: Number(query.get('page')) || 1,
		inLibraryOnly: query.get('inLibraryOnly') === 'true' ? true : undefined,
		status: query.get('status') as
			| 'watching'
			| 'completed'
			| 'planning'
			| 'dropped'
			| 'paused'
			| undefined,
	};
}

export interface FilterActions {
	updateSearch: (search: string) => void;
	toggleGenre: (genre: string) => void;
	toggleSeason: (season: FilterState['seasons'][number]) => void;
	toggleYear: (year: number) => void;
	setYearRange: (min?: number, max?: number) => void;
	setGenresMode: (mode: FilterState['genresMode']) => void;
	setSort: (sortBy: FilterState['sortBy'], sortOrder: FilterState['sortOrder']) => void;
	setPage: (page: number) => void;
	setItemsPerPage: (count: number) => void;
	setStatus: (status: FilterState['status']) => void;
	clearSearch: () => void;
	clearAll: () => void;
}

export function getTotalFilters(filters: FilterState): number {
	return (
		filters.genres.length +
		filters.seasons.length +
		filters.years.length +
		(filters.yearMin ? 1 : 0) +
		(filters.yearMax ? 1 : 0) +
		(filters.search ? 1 : 0)
	);
}
