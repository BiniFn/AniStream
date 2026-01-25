import { redirect } from '@sveltejs/kit';
import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	// Temporary redirect to 404 page for maintenance
	redirect(302, '/404');

	const [latestRelease, allReleases] = await Promise.all([
		apiClient.GET('/desktop/releases/latest', { fetch }),
		apiClient.GET('/desktop/releases', { fetch }),
	]);

	return {
		latestRelease: latestRelease.data,
		allReleases: allReleases.data || [],
	};
};
