import { apiClient } from '$lib/api/client';
import { filtersToApiQuery, searchParamsToFilters } from '$lib/utils/filters';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url }) => {
	const filters = searchParamsToFilters(url.searchParams);

	const apiQuery = filtersToApiQuery(filters);

	try {
		const [listings, genresList] = await Promise.all([
			apiClient.GET('/anime/listings', {
				fetch,
				params: {
					query: apiQuery,
				},
			}),
			apiClient.GET('/anime/listings/genres', { fetch }),
		]);

		return {
			listings: listings.data || {
				items: [],
				pageInfo: {
					currentPage: 1,
					totalPages: 1,
					hasNextPage: false,
					hasPrevPage: false,
				},
			},
			genres: genresList.data || [],
			initialFilters: filters,
		};
	} catch {
		return {
			listings: {
				items: [],
				pageInfo: {
					currentPage: 1,
					totalPages: 1,
					hasNextPage: false,
					hasPrevPage: false,
				},
			},
			genres: [],
			initialFilters: filters,
		};
	}
};
