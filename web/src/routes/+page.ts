import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const [seasonalAnime, recentlyUpdatedAnime] = await Promise.all([
		apiClient.GET('/anime/listings/seasonal', { fetch }),
		apiClient.GET('/anime/listings/recently-updated', { fetch }),
	]);

	return {
		seasonalAnime: seasonalAnime.data ?? [],
		recentlyUpdatedAnime: recentlyUpdatedAnime.data ?? [],
	};
};
