export const isMobile = () => {
	return window.innerWidth < 640;
};

export const isTablet = () => {
	return window.innerWidth >= 640 && window.innerWidth < 1024;
};

export const isPc = () => {
	return window.innerWidth >= 1024 && window.innerWidth <= 1980;
};

export const matchesBreakpoint = (min, max = Infinity) => {
	return window.innerWidth >= min && window.innerWidth <= max;
};
