import { apiClient } from '$lib/api/client';
import { redirectToErrorPage } from '$lib/errors';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	try {
		const character = await apiClient.GET('/characters/{id}', {
			fetch,
			params: { path: { id: parseInt(params.id) } },
		});

		if (character.error) {
			console.error('Failed to fetch character:', character.error);
			redirectToErrorPage('character_not_found_or_unavailable', `/characters/${params.id}`);
		}

		return {
			character: character.data as NonNullable<typeof character.data>,
		};
	} catch (error) {
		console.error('Error loading character:', error);
		redirectToErrorPage('character_not_found_or_unavailable', `/characters/${params.id}`);
	}
};
