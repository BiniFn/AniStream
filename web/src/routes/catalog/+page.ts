import { apiClient } from '$lib/api/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url }) => {
	const page = Number(url.searchParams.get('page') || '1');
	const itemsPerPageParam = url.searchParams.get('itemsPerPage');
	const itemsPerPage = itemsPerPageParam ? Number(itemsPerPageParam) : 24;
	const search = url.searchParams.get('search') || undefined;
	const genres = url.searchParams.getAll('genres');
	const genresMode = (url.searchParams.get('genresMode') as 'any' | 'all' | null) || undefined;
	const seasons = url.searchParams.getAll('seasons') as Array<
		'winter' | 'spring' | 'summer' | 'fall' | 'unknown'
	>;
	const years = url.searchParams.getAll('years').map(Number).filter(Boolean);
	const yearMinStr = url.searchParams.get('yearMin');
	const yearMaxStr = url.searchParams.get('yearMax');
	const sortBy =
		(url.searchParams.get('sortBy') as
			| 'ename'
			| 'jname'
			| 'season'
			| 'year'
			| 'relevance'
			| 'updated_at'
			| null) || undefined;
	const sortOrder = (url.searchParams.get('sortOrder') as 'asc' | 'desc' | null) || undefined;

	const [listings, genresList] = await Promise.all([
		apiClient.GET('/anime/listings', {
			fetch,
			params: {
				query: {
					page,
					itemsPerPage,
					search,
					genres: genres.length ? genres : undefined,
					genresMode,
					seasons: seasons.length ? seasons : undefined,
					years: years.length ? years : undefined,
					yearMin: yearMinStr ? Number(yearMinStr) : undefined,
					yearMax: yearMaxStr ? Number(yearMaxStr) : undefined,
					sortBy,
					sortOrder,
				},
			},
		}),
		apiClient.GET('/anime/listings/genres', { fetch }),
	]);

	return {
		listings: listings.data!,
		genres: genresList.data || [],
		initialQuery: {
			page,
			itemsPerPage,
			search: search ?? '',
			genres,
			genresMode: genresMode ?? 'any',
			seasons,
			years,
			yearMin: yearMinStr ? Number(yearMinStr) : undefined,
			yearMax: yearMaxStr ? Number(yearMaxStr) : undefined,
			sortBy: sortBy ?? 'relevance',
			sortOrder: sortOrder ?? 'desc',
		},
	};
};
