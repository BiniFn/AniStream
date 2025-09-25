import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return {
		stream: apiClient.GET('/anime/listings/random', { fetch }).then((res) => res.data),
	};
};
