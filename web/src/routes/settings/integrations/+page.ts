import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const response = await apiClient.GET('/auth/providers', { fetch });
		return {
			oauthProviders: response.data || [],
		};
	} catch {
		return {
			oauthProviders: [],
		};
	}
};
