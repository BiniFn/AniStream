import { goto } from '$app/navigation';
import type { FilterState } from './filters';
import { filtersToQuery, defaultFilters } from './filters';

export function createFilterActions(
	getFilters: () => FilterState,
	updateFilters: (newFilters: FilterState) => void,
) {
	return {
		updateSearch: (search: string) => {
			const filters = getFilters();
			updateFilters({ ...filters, search, page: 1 });
		},

		toggleGenre: (genre: string) => {
			const filters = getFilters();
			const newGenres = filters.genres.includes(genre)
				? filters.genres.filter((g) => g !== genre)
				: [...filters.genres, genre];
			updateFilters({ ...filters, genres: newGenres, page: 1 });
		},

		setGenresMode: (mode: FilterState['genresMode']) => {
			const filters = getFilters();
			updateFilters({ ...filters, genresMode: mode });
		},

		toggleSeason: (season: FilterState['seasons'][number]) => {
			const filters = getFilters();
			const newSeasons = filters.seasons.includes(season)
				? filters.seasons.filter((s) => s !== season)
				: [...filters.seasons, season];
			updateFilters({ ...filters, seasons: newSeasons, page: 1 });
		},

		toggleYear: (year: number) => {
			const filters = getFilters();
			const newYears = filters.years.includes(year)
				? filters.years.filter((y) => y !== year)
				: [...filters.years, year];
			updateFilters({ ...filters, years: newYears, page: 1 });
		},

		setYearRange: (yearMin?: number, yearMax?: number) => {
			const filters = getFilters();
			updateFilters({ ...filters, yearMin, yearMax, page: 1 });
		},

		setSort: (sortBy: FilterState['sortBy'], sortOrder: FilterState['sortOrder']) => {
			const filters = getFilters();
			updateFilters({ ...filters, sortBy, sortOrder });
		},

		setPage: (page: number) => {
			const filters = getFilters();
			updateFilters({ ...filters, page });
		},

		setItemsPerPage: (itemsPerPage: number) => {
			const filters = getFilters();
			updateFilters({ ...filters, itemsPerPage, page: 1 });
		},

		setStatus: (status: FilterState['status']) => {
			const filters = getFilters();
			updateFilters({ ...filters, status, page: 1 });
		},

		clearSearch: () => {
			const filters = getFilters();
			updateFilters({ ...filters, search: '', page: 1 });
		},

		clearAll: () => {
			updateFilters({ ...defaultFilters });
		},
	};
}

export async function updateUrlWithFilters(filters: FilterState) {
	const params = filtersToQuery(filters);
	await goto(`?${params.toString()}`, { replaceState: true });
}
