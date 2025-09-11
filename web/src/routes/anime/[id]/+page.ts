import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	const [anime, banner, trailer, episodes] = await Promise.allSettled([
		apiClient.GET('/anime/{id}', { fetch, params: { path: params } }),
		apiClient.GET('/anime/{id}/banner', { fetch, params: { path: params } }),
		apiClient.GET('/anime/{id}/trailer', { fetch, params: { path: params } }),
		apiClient.GET('/anime/{id}/episodes', { fetch, params: { path: params } }),
	]);

	return {
		anime: anime.status === 'fulfilled' ? anime.value : null,
		banner: banner.status === 'fulfilled' ? banner.value : null,
		trailer: trailer.status === 'fulfilled' ? trailer.value : null,
		episodes: episodes.status === 'fulfilled' ? episodes.value : null,
	};
};
