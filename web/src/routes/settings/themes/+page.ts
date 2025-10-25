import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const response = await apiClient.GET('/themes', { fetch });
		return {
			themes: response.data || [],
		};
	} catch (error) {
		console.error('Failed to load themes:', error);
		return {
			themes: [],
		};
	}
};
