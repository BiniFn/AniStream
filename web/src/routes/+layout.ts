import { apiClient } from '$lib/api/client';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ fetch }) => {
	const seasonalAnime = await apiClient.GET('/anime/listings/seasonal', { fetch });

	return { seasonalAnime };
};
