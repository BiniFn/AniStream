import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent }) => {
	const { user } = await parent();

	const [trending, popular, recentlyUpdated, seasonal, continueWatching, planning] =
		await Promise.allSettled([
			apiClient.GET('/anime/listings/trending', { fetch }),
			apiClient.GET('/anime/listings/popular', { fetch }),
			apiClient.GET('/anime/listings/recently-updated', {
				fetch,
				params: { query: { page: 1, itemsPerPage: 6 } },
			}),
			apiClient.GET('/anime/listings/seasonal', { fetch }),
			user
				? apiClient.GET('/library/continue-watching', {
						fetch,
						params: { query: { page: 1, itemsPerPage: 6 } },
					})
				: Promise.resolve({ data: { items: [] } }),
			user
				? apiClient.GET('/library/planning', {
						fetch,
						params: { query: { page: 1, itemsPerPage: 6 } },
					})
				: Promise.resolve({ data: { items: [] } }),
		]);

	const featuredAnime = trending.status === 'fulfilled' ? trending.value?.data?.[0] : null;
	const metadata = featuredAnime
		? await apiClient.GET('/anime/{id}', {
				fetch,
				params: { path: { id: featuredAnime?.id || '' } },
			})
		: null;

	return {
		trending: trending.status === 'fulfilled' ? trending.value.data || [] : [],
		popular: popular.status === 'fulfilled' ? popular.value.data || [] : [],
		recentlyUpdated:
			recentlyUpdated.status === 'fulfilled' ? recentlyUpdated.value.data?.items || [] : [],
		seasonal: seasonal.status === 'fulfilled' ? seasonal.value.data || [] : [],
		featuredAnime: metadata?.data,
		user,
		continueWatching:
			continueWatching.status === 'fulfilled' ? continueWatching.value.data?.items || [] : [],
		planning: planning.status === 'fulfilled' ? planning.value.data?.items || [] : [],
	};
};
