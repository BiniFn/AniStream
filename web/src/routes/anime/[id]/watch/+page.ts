import { apiClient } from '$lib/api/client';
import { redirectToErrorPage } from '$lib/errors';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, url }) => {
	const episodeNumber = Number(url.searchParams.get('ep')) || 1;

	const [anime, episodes, library] = await Promise.allSettled([
		apiClient.GET('/anime/{id}', { fetch, params: { path: params } }),
		apiClient.GET('/anime/{id}/episodes', { fetch, params: { path: params } }),
		apiClient
			.GET('/library/{animeID}', { fetch, params: { path: { animeID: params.id } } })
			.catch(() => null),
	]);

	const animeData = anime.status === 'fulfilled' ? anime.value?.data : null;
	const episodesData = episodes.status === 'fulfilled' ? episodes.value?.data : null;
	const libraryData = library.status === 'fulfilled' ? library.value?.data : null;

	if (!animeData || !episodesData) {
		redirectToErrorPage('anime_not_found_or_unavailable', `/anime/${params.id}`);
	}

	const currentEpisode = episodesData?.find((ep) => ep.number === episodeNumber);

	if (!currentEpisode) {
		redirectToErrorPage(
			'episode_not_found_or_unavailable',
			`/anime/${params.id}/watch?ep=${episodeNumber}`,
		);
	}

	const episodeServers = await apiClient.GET('/anime/{id}/episodes/{episodeID}/servers', {
		fetch,
		params: { path: { id: params.id, episodeID: currentEpisode!.id } },
	});

	return {
		anime: animeData as NonNullable<typeof animeData>,
		episodes: episodesData as NonNullable<typeof episodesData>,
		currentEpisode: currentEpisode as NonNullable<typeof currentEpisode>,
		episodeNumber,
		servers: episodeServers.data || [],
		libraryEntry: libraryData,
	};
};
