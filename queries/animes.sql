-- name: GetCountOfAnimes :one
SELECT COUNT(*) AS count
FROM animes;
-- name: GetAllGenres :many
SELECT DISTINCT trim(unnested) AS genre
FROM animes,
  unnest(string_to_array(genre, ',')) AS unnested
ORDER BY genre;
-- name: GetRecentlyUpdatedAnimes :many
SELECT *
FROM animes
WHERE animes.mal_id IS NOT NULL
  OR animes.mal_id != 0
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;
-- name: GetRecentlyUpdatedAnimesCount :one
SELECT COUNT(*)
FROM animes
WHERE animes.mal_id IS NOT NULL
  OR animes.mal_id != 0;
-- name: GetAnimeByGenre :many
SELECT *
FROM animes
WHERE genre ILIKE '%' || sqlc.arg(genre) || '%'
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;
-- name: GetAnimeByGenreCount :one
SELECT COUNT(*)
FROM animes
WHERE genre ILIKE '%' || sqlc.arg(genre) || '%';
-- name: GetRandomAnime :one
SELECT *
FROM animes
ORDER BY RANDOM()
LIMIT 1;
-- name: GetRandomAnimeByGenre :one
SELECT *
FROM animes
WHERE genre ILIKE '%' || sqlc.arg(genre) || '%'
ORDER BY RANDOM()
LIMIT 1;
-- name: GetAnimeById :one
SELECT *
FROM animes
WHERE id = $1;
-- name: GetAnimeByMalId :one
SELECT *
FROM animes
WHERE mal_id = $1;
-- name: GetAnimeByAnilistId :one
SELECT *
FROM animes
WHERE anilist_id = $1;
-- name: GetAnimeByHiAnimeId :one
SELECT *
FROM animes
WHERE hi_anime_id = $1;
-- name: GetAnimeMetadataByMalId :one
SELECT *
FROM anime_metadata
WHERE mal_id = $1;
-- name: GetAnimesByIds :many
SELECT *
FROM animes
WHERE id = ANY(sqlc.arg(ids)::text [])
ORDER BY updated_at DESC;
-- name: GetAnimesByMalIds :many
SELECT *
FROM animes
WHERE mal_id = ANY(sqlc.arg(mal_ids)::int [])
ORDER BY updated_at DESC;
-- name: GetAnimesByAnilistIds :many
SELECT *
FROM animes
WHERE anilist_id = ANY(sqlc.arg(anilist_ids)::int [])
ORDER BY updated_at DESC;
-- name: GetAnimesByHiAnimeIds :many
SELECT *
FROM animes
WHERE hi_anime_id = ANY(sqlc.arg(hi_anime_ids)::text [])
ORDER BY updated_at DESC;
-- name: SearchAnimes :many
SELECT animes.*,
  ts_rank(
    animes.search_vector,
    plainto_tsquery(sqlc.arg(query))
  ) AS query_rank
FROM animes
WHERE (
    sqlc.arg(query) = ''
    OR sqlc.arg(query) IS NULL
    OR ename % sqlc.arg(query)
    OR jname % sqlc.arg(query)
    OR search_vector @@ plainto_tsquery('english', sqlc.arg(query))
  )
  AND (
    sqlc.arg(genre) = ''
    OR sqlc.arg(genre) IS NULL
    OR genre ILIKE '%' || sqlc.arg(genre) || '%'
  )
  AND animes.mal_id IS NOT NULL
ORDER BY query_rank DESC
LIMIT $1 OFFSET $2;
-- name: SearchAnimesCount :one
SELECT COUNT(*)
FROM animes
WHERE (
    sqlc.arg(query) = ''
    OR sqlc.arg(query) IS NULL
    OR ename % sqlc.arg(query)
    OR jname % sqlc.arg(query)
    OR search_vector @@ plainto_tsquery('english', sqlc.arg(query))
  )
  AND (
    sqlc.arg(genre) = ''
    OR sqlc.arg(genre) IS NULL
    OR genre ILIKE '%' || sqlc.arg(genre) || '%'
  )
  AND animes.mal_id IS NOT NULL;
-- name: InsertAnime :exec
INSERT INTO animes (
    ename,
    jname,
    image_url,
    genre,
    hi_anime_id,
    mal_id,
    anilist_id,
    last_episode,
    created_at,
    updated_at
  )
VALUES (
    sqlc.arg(ename),
    sqlc.arg(jname),
    sqlc.arg(image_url),
    sqlc.arg(genre),
    sqlc.arg(hi_anime_id),
    sqlc.arg(mal_id),
    sqlc.arg(anilist_id),
    sqlc.arg(last_episode),
    COALESCE(sqlc.arg(created_at), NOW()),
    COALESCE(sqlc.arg(updated_at), NOW())
  )
RETURNING *;
-- name: InsertMultipleAnimes :copyfrom
INSERT INTO animes (
    ename,
    jname,
    image_url,
    genre,
    hi_anime_id,
    mal_id,
    anilist_id,
    last_episode
  )
VALUES (
    sqlc.arg(ename),
    sqlc.arg(jname),
    sqlc.arg(image_url),
    sqlc.arg(genre),
    sqlc.arg(hi_anime_id),
    sqlc.arg(mal_id),
    sqlc.arg(anilist_id),
    sqlc.arg(last_episode)
  );
-- name: UpdateAnime :exec
UPDATE animes
SET ename = sqlc.arg(ename),
  jname = sqlc.arg(jname),
  image_url = sqlc.arg(image_url),
  genre = sqlc.arg(genre),
  hi_anime_id = sqlc.arg(hi_anime_id),
  mal_id = sqlc.arg(mal_id),
  anilist_id = sqlc.arg(anilist_id),
  last_episode = sqlc.arg(last_episode),
  updated_at = COALESCE(sqlc.arg(updated_at), NOW())
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: UpsertAnimeMetadata :exec
INSERT INTO anime_metadata (
    mal_id,
    description,
    main_picture_url,
    media_type,
    rating,
    airing_status,
    avg_episode_duration,
    total_episodes,
    studio,
    rank,
    mean,
    scoringUsers,
    popularity,
    airing_start_date,
    airing_end_date,
    source,
    trailer_embed_url,
    season_year,
    season
  )
VALUES (
    sqlc.arg(mal_id),
    sqlc.arg(description),
    sqlc.arg(main_picture_url),
    sqlc.arg(media_type),
    sqlc.arg(rating),
    sqlc.arg(airing_status),
    sqlc.arg(avg_episode_duration),
    sqlc.arg(total_episodes),
    sqlc.arg(studio),
    sqlc.arg(rank),
    sqlc.arg(mean),
    sqlc.arg(scoringUsers),
    sqlc.arg(popularity),
    sqlc.arg(airing_start_date),
    sqlc.arg(airing_end_date),
    sqlc.arg(source),
    sqlc.arg(trailer_embed_url),
    sqlc.arg(season_year),
    sqlc.arg(season)
  ) ON CONFLICT (mal_id) DO
UPDATE
SET description = EXCLUDED.description,
  main_picture_url = EXCLUDED.main_picture_url,
  media_type = EXCLUDED.media_type,
  rating = EXCLUDED.rating,
  airing_status = EXCLUDED.airing_status,
  avg_episode_duration = EXCLUDED.avg_episode_duration,
  total_episodes = EXCLUDED.total_episodes,
  studio = EXCLUDED.studio,
  rank = EXCLUDED.rank,
  mean = EXCLUDED.mean,
  scoringUsers = EXCLUDED.scoringUsers,
  popularity = EXCLUDED.popularity,
  airing_start_date = EXCLUDED.airing_start_date,
  airing_end_date = EXCLUDED.airing_end_date,
  source = EXCLUDED.source,
  trailer_embed_url = EXCLUDED.trailer_embed_url,
  season_year = EXCLUDED.season_year,
  season = EXCLUDED.season,
  updated_at = NOW()
RETURNING *;
-- name: InsertMultipleAnimeMetadatas :copyfrom
INSERT INTO anime_metadata (
    mal_id,
    description,
    main_picture_url,
    media_type,
    rating,
    airing_status,
    avg_episode_duration,
    total_episodes,
    studio,
    rank,
    mean,
    scoringUsers,
    popularity,
    airing_start_date,
    airing_end_date,
    source,
    trailer_embed_url,
    season_year,
    season
  )
VALUES (
    sqlc.arg(mal_id),
    sqlc.arg(description),
    sqlc.arg(main_picture_url),
    sqlc.arg(media_type),
    sqlc.arg(rating),
    sqlc.arg(airing_status),
    sqlc.arg(avg_episode_duration),
    sqlc.arg(total_episodes),
    sqlc.arg(studio),
    sqlc.arg(rank),
    sqlc.arg(mean),
    sqlc.arg(scoringUsers),
    sqlc.arg(popularity),
    sqlc.arg(airing_start_date),
    sqlc.arg(airing_end_date),
    sqlc.arg(source),
    sqlc.arg(trailer_embed_url),
    sqlc.arg(season_year),
    sqlc.arg(season)
  );
-- name: UpdateAnimeMetadataTrailer :exec
UPDATE anime_metadata
SET trailer_embed_url = sqlc.arg(trailer_embed_url),
  updated_at = NOW()
WHERE mal_id = sqlc.arg(mal_id)
RETURNING *;