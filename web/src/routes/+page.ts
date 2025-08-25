import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const seasonalAnime = await apiClient.GET('/anime/listings/seasonal', {
		fetch,
	});

	return { seasonalAnime };
};
