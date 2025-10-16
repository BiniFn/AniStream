import { getContext, setContext } from 'svelte';
import type { Attachment } from 'svelte/attachments';

export class LayoutState {
	navbarHeight = $state(68);
	headerHeight = $state(0);
	totalHeight = $derived(this.navbarHeight + this.headerHeight);

	setHeight = (type: 'navbar' | 'header'): Attachment => {
		return (el) => {
			const updateHeight = () => {
				if (el) {
					const height = el.getBoundingClientRect().height;
					if (type === 'navbar') {
						this.navbarHeight = height;
					} else if (type === 'header') {
						this.headerHeight = height;
					}
				}
			};

			const resizeObserver = new ResizeObserver(() => {
				updateHeight();
			});

			resizeObserver.observe(el);
			updateHeight();

			window.addEventListener('resize', updateHeight);
			return () => {
				window.removeEventListener('resize', updateHeight);
				if (resizeObserver && el) {
					resizeObserver.unobserve(el);
					resizeObserver.disconnect();
				}
			};
		};
	};
}

const LAYOUT_STATE = Symbol('LAYOUT_STATE');

export const setLayoutStateContext = (key = LAYOUT_STATE) => {
	const layoutState = new LayoutState();
	return setContext(key, layoutState);
};

export const getLayoutStateContext = (key = LAYOUT_STATE) => {
	const layoutState = getContext<LayoutState>(key);
	if (!layoutState) {
		throw new Error('LayoutState context not found');
	}
	return layoutState;
};
