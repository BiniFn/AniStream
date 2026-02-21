import type { Handle, HandleFetch } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

export const handle: Handle = async ({ event, resolve }) => {
	return resolve(event, {
		filterSerializedResponseHeaders: (name) => name === 'content-type' || name === 'content-length',
	});
};

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	const apiUrl = env.PUBLIC_API_URL?.trim();

	if (apiUrl && request.url.startsWith(apiUrl)) {
		const cookies = event.request.headers.get('cookie');
		if (cookies) {
			request = new Request(request, {
				headers: new Headers(request.headers),
			});
			request.headers.set('cookie', cookies);
		}
	}

	return fetch(request);
};
