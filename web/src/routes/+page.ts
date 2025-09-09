import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const [trending, popular, recentlyUpdated, seasonal] = await Promise.all([
		apiClient.GET('/anime/listings/trending', { fetch }),
		apiClient.GET('/anime/listings/popular', { fetch }),
		apiClient.GET('/anime/listings/recently-updated', {
			fetch,
			params: { query: { page: 1, itemsPerPage: 6 } },
		}),
		apiClient.GET('/anime/listings/seasonal', { fetch }),
	]);

	const featuredAnime = trending?.data?.[0];
	const metadata = featuredAnime
		? await apiClient.GET('/anime/{id}', {
				fetch,
				params: { path: { id: featuredAnime?.id || '' } },
			})
		: null;
	return {
		trending: trending.data || [],
		popular: popular.data || [],
		recentlyUpdated: recentlyUpdated.data?.items || [],
		seasonal: seasonal.data || [],
		featuredAnime: metadata?.data,
	};
};
