import { z } from "zod";
import { createPaginationSchema } from "./common";

export const AnimeMetadataSchema = z.object({
  malId: z.number(),
  description: z.string(),
  mainPictureUrl: z.string(),
  mediaType: z.string(),
  rating: z.string(),
  airingStatus: z.string(),
  avgEpisodeDuration: z.number(),
  totalEpisodes: z.number(),
  studio: z.string(),
  rank: z.number(),
  mean: z.number(),
  scoringUsers: z.number(),
  popularity: z.number(),
  airingStartDate: z.string(),
  airingEndDate: z.string(),
  source: z.string(),
  seasonYear: z.number(),
  season: z.string(),
  trailerEmbedUrl: z.string(),
});

export const AnimeSchema = z.object({
  id: z.string(),
  ename: z.string().nullable(),
  jname: z.string().nullable(),
  imageUrl: z.string(),
  genre: z.string(),
  malId: z.number().nullable(),
  anilistId: z.number().nullable(),
  lastEpisode: z.number().nullable(),
  metadata: AnimeMetadataSchema.nullable().optional(),
});

export const PaginatedAnimeSchema = createPaginationSchema(AnimeSchema);

export const TrailerSchema = z.object({
  trailer: z.string(),
});

export const EpisodeSchema = z.object({
  id: z.string(),
  title: z.string(),
  number: z.number(),
  isFiller: z.boolean(),
});

export const EpisodeSourceSchema = z.object({
  url: z.string(),
  rawUrl: z.string(),
});

export const SegmentSchema = z.object({
  start: z.number(),
  end: z.number(),
});

export const TrackSchema = z.object({
  url: z.string(),
  raw: z.string(),
  kind: z.string(),
  label: z.string().optional(),
  default: z.boolean().optional(),
});

export const StreamingMetadataSchema = z.object({
  intro: SegmentSchema,
  outro: SegmentSchema,
  tracks: z.array(TrackSchema),
});

export const SeasonalAnimeSchema = z.object({
  id: z.string(),
  bannerImageUrl: z.string(),
  description: z.string(),
  startDate: z.number(),
  type: z.string(),
  episodes: z.number(),
  anime: AnimeSchema,
});

export const RelationsSchema = z.object({
  watchOrder: z.array(AnimeSchema),
  related: z.array(AnimeSchema),
});

export const BannerSchema = z.object({
  url: z.string(),
});
