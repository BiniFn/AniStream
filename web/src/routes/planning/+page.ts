import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('app:library');
	const page = parseInt(url.searchParams.get('page') || '1');

	const listings = await apiClient.GET('/library/planning', {
		fetch,
		params: {
			query: {
				page,
				itemsPerPage: 24,
			},
		},
	});

	return {
		listings: listings.data || {
			items: [],
			pageInfo: { currentPage: 1, totalPages: 1, hasNextPage: false, hasPrevPage: false },
		},
	};
};
