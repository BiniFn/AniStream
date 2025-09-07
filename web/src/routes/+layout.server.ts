import { apiClient } from '$lib/api/client';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
	try {
		const [user, settings] = await Promise.all([
			apiClient.GET('/auth/me', { fetch }),
			apiClient.GET('/settings', { fetch }),
		]);

		return {
			user: user.data || null,
			settings: settings.data || null,
			isLoggedIn: !!user.data,
		};
	} catch {
		return {
			user: null,
			settings: null,
			isLoggedIn: false,
		};
	}
};
