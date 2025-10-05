import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';
import { redirectToErrorPage } from '$lib/errors';

export const load: PageLoad = async ({ fetch, params }) => {
	const [anime, banner, trailer, episodes, franchise, libraryStatus, characters] =
		await Promise.allSettled([
			apiClient.GET('/anime/{id}', { fetch, params: { path: params } }),
			apiClient.GET('/anime/{id}/banner', { fetch, params: { path: params } }),
			apiClient.GET('/anime/{id}/trailer', { fetch, params: { path: params } }),
			apiClient.GET('/anime/{id}/episodes', { fetch, params: { path: params } }),
			apiClient.GET('/anime/{id}/franchise', { fetch, params: { path: params } }),
			apiClient.GET('/library/{animeID}', { fetch, params: { path: { animeID: params.id } } }),
			apiClient.GET('/anime/{id}/characters', { fetch, params: { path: params } }),
		]);

	const isAnyRejected = [
		anime,
		banner,
		trailer,
		episodes,
		franchise,
		libraryStatus,
		characters,
	].some((result) => result.status === 'rejected');
	if (isAnyRejected) {
		console.error('One or more API requests failed:', {
			anime,
			banner,
			trailer,
			episodes,
			franchise,
			characters,
		});
	}

	const isAnimeError = anime.status === 'rejected';
	const animeData = anime.status === 'fulfilled' ? anime.value?.data : null;
	if (isAnimeError || !animeData) {
		redirectToErrorPage('anime_not_found_or_unavailable');
	}

	const ratings: Record<string, string> = {
		g: 'G - All Ages',
		pg: 'PG - Children',
		pg_13: 'PG-13 - Teens 13+',
		r: 'R - 17+ (violence & profanity)',
		r_plus: 'R+ - Mild Nudity',
		rx: 'Rx - Hentai',
		unknown: 'Unknown Rating',
	};

	return {
		anime: animeData as NonNullable<typeof animeData>,
		ratingLabel: ratings[animeData?.metadata?.rating ?? 'unknown'] || 'Unknown Rating',
		banner: banner.status === 'fulfilled' ? banner.value : null,
		trailer: trailer.status === 'fulfilled' ? trailer.value : null,
		episodes: episodes.status === 'fulfilled' ? episodes.value : null,
		franchise: franchise.status === 'fulfilled' ? franchise.value : null,
		libraryStatus: libraryStatus.status === 'fulfilled' ? libraryStatus.value : null,
		characters: characters.status === 'fulfilled' ? characters.value : null,
	};
};
