import { apiClient } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import type { PageLoad } from './$types';

export const prerender = true;

type GenrePreview = components['schemas']['models.GenrePreview'];

export const load: PageLoad = async ({ fetch }) => {
  const previews = await apiClient.GET('/anime/listings/genres/previews', { fetch });
  const genres = (previews.data || []).filter((g) => g.name?.toLowerCase() !== 'unknown');
  return { genres } satisfies { genres: GenrePreview[] };
};
