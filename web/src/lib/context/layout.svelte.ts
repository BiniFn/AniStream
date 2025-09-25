type LayoutState = {
	navbarHeight: number;
	headerHeight: number;
};

export const layoutState = $state<LayoutState>({
	navbarHeight: 68,
	headerHeight: 0,
});

export const setNavbarHeight = (height: number) => {
	layoutState.navbarHeight = height;
};

export const setHeaderHeight = (height: number) => {
	layoutState.headerHeight = height;
};

