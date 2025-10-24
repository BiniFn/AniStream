import { goto } from '$app/navigation';
import {
	type FilterState,
	defaultFilters,
	getTotalFilters,
	filtersToSearchParams,
	type FilterActions,
} from './filters';

export class FilterManager implements FilterActions {
	filters = $state<FilterState>({ ...defaultFilters });
	isLoading = $state(false);
	totalFilters = $derived.by(() => getTotalFilters(this.filters));

	viewMode = $state<'grid' | 'list'>('grid');
	showMobileFilters = $state(false);

	constructor(initial?: FilterState) {
		if (initial) this.filters = initial;
	}

	private async _updateFilters(newFilters: FilterState) {
		this.isLoading = true;
		this.filters = newFilters;
		const params = filtersToSearchParams(this.filters);
		await goto(`?${params.toString()}`, { replaceState: true });
		this.isLoading = false;
	}

	updateSearch = async (search: string) => {
		await this._updateFilters({ ...this.filters, search, page: 1 });
	};

	clearSearch = async () => {
		await this._updateFilters({ ...this.filters, search: '', page: 1 });
	};

	toggleGenre = async (genre: string) => {
		const newGenres = this.filters.genres.includes(genre)
			? this.filters.genres.filter((g) => g !== genre)
			: [...this.filters.genres, genre];
		await this._updateFilters({ ...this.filters, genres: newGenres, page: 1 });
	};

	setGenresMode = async (mode: FilterState['genresMode']) => {
		await this._updateFilters({ ...this.filters, genresMode: mode, page: 1 });
	};

	toggleSeason = async (season: FilterState['seasons'][number]) => {
		const newSeasons = this.filters.seasons.includes(season)
			? this.filters.seasons.filter((s) => s !== season)
			: [...this.filters.seasons, season];
		await this._updateFilters({ ...this.filters, seasons: newSeasons, page: 1 });
	};

	toggleYear = async (year: number) => {
		const newYears = this.filters.years.includes(year)
			? this.filters.years.filter((y) => y !== year)
			: [...this.filters.years, year];
		await this._updateFilters({ ...this.filters, years: newYears, page: 1 });
	};

	setYearRange = async (min?: number, max?: number) => {
		await this._updateFilters({ ...this.filters, yearMin: min, yearMax: max, page: 1 });
	};

	setSort = async (sortBy: FilterState['sortBy'], sortOrder: FilterState['sortOrder']) => {
		await this._updateFilters({ ...this.filters, sortBy, sortOrder, page: 1 });
	};

	setPage = async (page: number) => {
		await this._updateFilters({ ...this.filters, page });
	};

	setItemsPerPage = async (count: number) => {
		await this._updateFilters({ ...this.filters, itemsPerPage: count, page: 1 });
	};

	setStatus = async (status: FilterState['status']) => {
		await this._updateFilters({ ...this.filters, status, page: 1 });
	};

	clearAll = async () => {
		await this._updateFilters({ ...defaultFilters });
	};

	handleViewModeChange = (mode: 'grid' | 'list') => {
		this.viewMode = mode;
	};

	setMobileFiltersVisibility = (visible: boolean) => {
		this.showMobileFilters = visible;
	};
}
