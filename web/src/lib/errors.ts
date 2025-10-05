import { redirect } from '@sveltejs/kit';

export const errors = [
	'anime_not_found_or_unavailable',
	'episode_not_found_or_unavailable',
	'character_not_found_or_unavailable',
	'person_not_found_or_unavailable',
] as const;

export type ErrorType = (typeof errors)[number];

export const errorMessages: Record<ErrorType, string> = {
	anime_not_found_or_unavailable: 'The requested anime was not found or is unavailable.',
	episode_not_found_or_unavailable: 'The requested episode was not found or is unavailable.',
	character_not_found_or_unavailable: 'The requested character was not found or is unavailable.',
	person_not_found_or_unavailable: 'The requested voice actor was not found or is unavailable.',
} as const;

export const redirectToErrorPage = (type: ErrorType) => {
	redirect(302, `/error?type=${type}`);
};
