import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';
import { redirectToErrorPage } from '$lib/errors';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends('app:library');
	const response = await apiClient.GET('/anime/{id}/full', {
		fetch,
		params: { path: params },
	});

	if (response.error || !response.data) {
		redirectToErrorPage('anime_not_found_or_unavailable', `/anime/${params.id}`);
	}

	const data = response.data!;

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
		anime: data.anime,
		ratingLabel: ratings[data.anime?.metadata?.rating ?? 'unknown'] || 'Unknown Rating',
		banner: data.banner,
		trailer: data.trailer,
		episodes: data.episodes,
		franchise: data.franchise,
		libraryStatus: data.libraryStatus,
		characters: data.characters,
	};
};
