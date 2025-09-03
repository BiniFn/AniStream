import type { components } from '$lib/api/openapi';

type UserResponse = components['schemas']['models.UserResponse'];

let user = $state<UserResponse | null>(null);

export function setUser(newUser: UserResponse | null) {
	user = newUser;
}

export function getUser() {
	return user;
}

export function isLoggedIn() {
	return user !== null;
}