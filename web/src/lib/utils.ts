import { clsx, type ClassValue } from 'clsx';
import { mount, type Component, type ComponentProps } from 'svelte';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export const convertComponentToHTML = (
	component: Parameters<typeof mount>[0],
	props: ComponentProps<Component>,
) => {
	const div = document.createElement('div');
	mount(component, { target: div, props });
	const string = div.innerHTML;
	div.remove();
	return string;
};
