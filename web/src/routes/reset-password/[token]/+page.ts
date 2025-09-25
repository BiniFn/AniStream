import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	try {
		const response = await apiClient.GET('/auth/u/{token}', {
			params: {
				path: { token: params.token },
			},
			fetch,
		});

		if (response.response.status === 200 && response.data) {
			return {
				token: params.token,
				user: response.data.user,
				expiresAt: response.data.expires_at,
				isValidToken: true,
			};
		}

		return {
			token: params.token,
			user: null,
			expiresAt: null,
			isValidToken: false,
		};
	} catch (error) {
		console.error('Token validation error:', error);
		return {
			token: params.token,
			user: null,
			isValidToken: false,
		};
	}
};

