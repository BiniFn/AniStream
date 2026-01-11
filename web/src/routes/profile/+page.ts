import { redirect } from '@sveltejs/kit';
import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent, depends }) => {
	depends('app:library');
	const { user } = await parent();

	if (!user) {
		redirect(302, '/login');
	}

	const response = await apiClient.GET('/library/stats', { fetch });

	if (response.error || !response.data) {
		// Fallback to zero stats if endpoint fails
		return {
			user,
			stats: {
				watching: 0,
				planning: 0,
				completed: 0,
			},
		};
	}

	return {
		user,
		stats: {
			watching: response.data.watching || 0,
			planning: response.data.planning || 0,
			completed: response.data.completed || 0,
		},
	};
};
