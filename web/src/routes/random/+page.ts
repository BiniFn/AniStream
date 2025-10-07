import { apiClient } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return {
		stream: apiClient
			.GET('/anime/listings/random', { fetch })
			.then((res) => res.data)
			.then(async (data) => {
				if (!data?.id) return new Promise(() => undefined);

				const res = await apiClient.GET('/anime/{id}', {
					params: { path: { id: data.id } },
					fetch,
				});

				return res.data;
			}) as Promise<components['schemas']['models.AnimeWithMetadataResponse'] | undefined>,
	};
};
