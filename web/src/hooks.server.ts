import type { HandleFetch } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	if (request.url.startsWith(PUBLIC_API_URL)) {
		request.headers.set('cookie', event.request.headers.get('cookie') || '');
	}

	return fetch(request);
};
