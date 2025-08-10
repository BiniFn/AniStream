export const useApiUrls = () => {
    const config = useRuntimeConfig();
    const apiBaseUrl = config.public.apiBaseUrl;

    return {
        anime: {
            listings: {
                recentlyUpdated: `${apiBaseUrl}/anime/listings/recently-updated`,
            },
        },
    };
};
