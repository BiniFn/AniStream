import { apiClient } from '$lib/api/client';
import { filtersToApiQuery, searchParamsToFilters } from '$lib/utils/filters';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url }) => {
	const filters = searchParamsToFilters(url.searchParams);

	const apiQuery = filtersToApiQuery(filters);

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
		listings: listings.data!,
		genres: genresList.data || [],
		initialFilters: filters,
	};
};
