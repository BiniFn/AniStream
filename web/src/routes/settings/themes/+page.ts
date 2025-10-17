import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const response = await apiClient.GET('/themes');
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

