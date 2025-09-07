// See https://svelte.dev/docs/kit/types#app.d.ts

import type { components } from '$lib/api/openapi';

// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		interface PageData {
			user: components['schemas']['models.UserResponse'] | null;
			settings: components['schemas']['models.SettingsResponse'] | null;
			isLoggedIn: boolean;
		}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
