import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/api/client';

export const load: PageLoad = async ({ fetch, parent, url }) => {
	const { user } = await parent();

	if (!user) {
		redirect(303, '/login');
	}

	const searchParams = url.searchParams;
	const status = (searchParams.get('status') || 'watching') as
		| 'watching'
		| 'planning'
		| 'completed'
		| 'dropped'
		| 'paused';
	const page = parseInt(searchParams.get('page') || '1', 10);
	const itemsPerPageParam = searchParams.get('itemsPerPage');
	const itemsPerPage = itemsPerPageParam ? Number(itemsPerPageParam) : 24;
	const search = searchParams.get('search') || undefined;
	const genres = searchParams.getAll('genres');
	const genresMode = (searchParams.get('genresMode') as 'any' | 'all' | null) || undefined;
	const seasons = searchParams.getAll('seasons') as Array<
		'winter' | 'spring' | 'summer' | 'fall' | 'unknown'
	>;
	const years = searchParams.getAll('years').map(Number).filter(Boolean);
	const yearMinStr = searchParams.get('yearMin');
	const yearMaxStr = searchParams.get('yearMax');
	const sortBy =
		(searchParams.get('sortBy') as
			| 'ename'
			| 'jname'
			| 'season'
			| 'year'
			| 'relevance'
			| 'updated_at'
			| 'anime_updated_at'
			| 'library_updated_at'
			| null) || undefined;
	const sortOrder = (searchParams.get('sortOrder') as 'asc' | 'desc' | null) || undefined;

	const [listings, genresList] = await Promise.all([
		apiClient.GET('/anime/listings', {
			fetch,
			params: {
				query: {
					inLibraryOnly: true,
					status,
					page,
					itemsPerPage,
					search,
					genres: genres.length ? genres : undefined,
					genresMode,
					seasons: seasons.length ? seasons : undefined,
					years: years.length ? years : undefined,
					yearMin: yearMinStr ? Number(yearMinStr) : undefined,
					yearMax: yearMaxStr ? Number(yearMaxStr) : undefined,
					sortBy: sortBy || 'library_updated_at',
					sortOrder: sortOrder || 'desc',
				},
			},
		}),
		apiClient.GET('/anime/listings/genres', { fetch }),
	]);

	if (listings.error) {
		console.error('Failed to fetch library:', listings.error);
		return {
			status,
			page,
			listings: {
				items: [],
				pageInfo: {
					currentPage: 1,
					totalPages: 1,
					hasNextPage: false,
					hasPrevPage: false,
				},
			},
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
				sortBy: sortBy ?? 'library_updated_at',
				sortOrder: sortOrder ?? 'desc',
			},
		};
	}

	return {
		status,
		page,
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
			sortBy: sortBy ?? 'library_updated_at',
			sortOrder: sortOrder ?? 'desc',
		},
	};
};
