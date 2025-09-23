import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/api/client';

export const load: PageLoad = async ({ fetch, parent, url }) => {
	const { user } = await parent();

	if (!user) {
		redirect(303, '/login');
	}

	const searchParams = url.searchParams;
	const status = (searchParams.get('status') || 'watching') as
		| 'watching'
		| 'planning'
		| 'completed'
		| 'dropped'
		| 'paused';
	const page = parseInt(searchParams.get('page') || '1', 10);
	const itemsPerPage = 24;

	const libraryData = await apiClient.GET('/library', {
		fetch,
		params: {
			query: {
				status,
				page,
				itemsPerPage,
			},
		},
	});

	if (libraryData.error) {
		console.error('Failed to fetch library:', libraryData.error);
		return {
			status,
			page,
			library: {
				items: [],
				pageInfo: {
					currentPage: 1,
					totalPages: 1,
					hasNextPage: false,
					hasPrevPage: false,
				},
			},
		};
	}

	return {
		status,
		page,
		library: libraryData.data,
	};
};
