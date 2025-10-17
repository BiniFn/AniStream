const THEME_FONT_URLS: Record<string, string> = {
	teal: 'https://fonts.googleapis.com/css2?family=Work+Sans:wght@400;600;700&display=swap',
	amber:
		'https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&family=Source+Serif+4:wght@400;600&family=JetBrains+Mono:wght@400;600&display=swap',
	catpuccin:
		'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;600;700&family=Fira+Code:wght@400;600&display=swap',
	cyberpunk:
		'https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;700&family=Fira+Code:wght@400;600&display=swap',
	ocean_breeze:
		'https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;700&family=Lora:wght@400;700&family=IBM+Plex+Mono:wght@400;700&display=swap',
	root: 'https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;600;700&family=Lora:wght@400;700&family=IBM+Plex+Mono:wght@400;700&display=swap',
};

export const getFontUrlsForTheme = (theme: string): string => {
	if (!theme || !(theme in THEME_FONT_URLS)) {
		theme = 'root';
	}

	return THEME_FONT_URLS[theme]!;
};
