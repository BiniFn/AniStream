import { redirect } from '@sveltejs/kit';
import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent }) => {
	const { user } = await parent();

	if (!user) {
		redirect(302, '/login');
	}

	const libraryStats = await apiClient.GET('/library', {
		fetch,
		params: {
			query: {
				status: 'watching',
				page: 1,
				itemsPerPage: 1,
			},
		},
	});

	const planningStats = await apiClient.GET('/library/planning', {
		fetch,
		params: {
			query: {
				page: 1,
				itemsPerPage: 1,
			},
		},
	});

	const completedStats = await apiClient.GET('/library', {
		fetch,
		params: {
			query: {
				status: 'completed',
				page: 1,
				itemsPerPage: 1,
			},
		},
	});

	return {
		user,
		stats: {
			watching: libraryStats.data?.pageInfo?.totalPages ? libraryStats.data.pageInfo.totalPages : 0,
			planning: planningStats.data?.pageInfo?.totalPages
				? planningStats.data.pageInfo.totalPages
				: 0,
			completed: completedStats.data?.pageInfo?.totalPages
				? completedStats.data.pageInfo.totalPages
				: 0,
		},
	};
};
