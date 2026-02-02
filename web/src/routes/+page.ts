import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends }) => {
	depends('app:library');

	const response = await apiClient.GET('/home', { fetch });

	if (response.error || !response.data) {
		console.error('Homepage data fetch error:', response.error);
		// Fallback to empty data if home endpoint fails
		return {
			trending: [],
			popular: [],
			recentlyUpdated: [],
			seasonal: [],
			featuredAnime: null,
			continueWatching: [],
			planning: [],
		};
	}

	const data = response.data;

	return {
		trending: data.trending || [],
		popular: data.popular || [],
		recentlyUpdated: data.recentlyUpdated || [],
		seasonal: data.seasonal || [],
		featuredAnime: data.featuredAnime || null,
		continueWatching: data.continueWatching || [],
		planning: data.planning || [],
	};
};
