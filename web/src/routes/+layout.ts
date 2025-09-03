import { apiClient } from '$lib/api/client';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ fetch }) => {
	const user = await apiClient
		.GET('/auth/me', { fetch })
		.catch(() => ({ data: null, error: null }));

	return {
		user: user.data,
	};
};
