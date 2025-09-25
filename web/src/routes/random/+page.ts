import { redirect } from '@sveltejs/kit';
import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const randomAnime = await apiClient.GET('/anime/listings/random', { fetch });

	if (randomAnime.data) {
		redirect(302, `/anime/${randomAnime.data.id}`);
	}

	redirect(302, '/catalog');
};

