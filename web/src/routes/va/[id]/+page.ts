import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';
import { redirectToErrorPage } from '$lib/errors';

export const load: PageLoad = async ({ fetch, params }) => {
	try {
		const person = await apiClient.GET('/characters/va/{id}', {
			fetch,
			params: { path: { id: parseInt(params.id) } },
		});

		if (person.error) {
			console.error('Failed to fetch person:', person.error);
			redirectToErrorPage('person_not_found_or_unavailable', `/va/${params.id}`);
		}

		return {
			person: person.data as NonNullable<typeof person.data>,
		};
	} catch (error) {
		console.error('Error loading person:', error);
		redirectToErrorPage('person_not_found_or_unavailable', `/va/${params.id}`);
	}
};
