import { redirect } from '@sveltejs/kit';
import { apiClient } from '$lib/api/client';
import { filtersToApiQuery, searchParamsToFilters } from '$lib/utils/filters';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent, url, depends }) => {
	depends('app:library');
	const { user } = await parent();

	if (!user) {
		redirect(303, '/login');
	}

	const filters = searchParamsToFilters(url.searchParams);
	const status = filters.status || 'watching';
	const apiQuery = filtersToApiQuery({
		...filters,
		inLibraryOnly: true,
		status,
		sortBy: filters.sortBy === 'relevance' ? 'library_updated_at' : filters.sortBy,
	});

	const [listings, genresList] = await Promise.all([
		apiClient.GET('/anime/listings', {
			fetch,
			params: {
				query: apiQuery,
			},
		}),
		apiClient.GET('/anime/listings/genres', { fetch }),
	]);

	if (listings.error) {
		console.error('Failed to fetch library:', listings.error);
		return {
			status,
			listings: {
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
	}

	return {
		status,
		listings: listings.data!,
		genres: genresList.data || [],
		initialFilters: filters,
	};
};
