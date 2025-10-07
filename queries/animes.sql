-- name: GetCountOfAnimes :one
SELECT
  COUNT(*) AS count
FROM
  animes;

-- name: GetAllGenres :many
SELECT DISTINCT
  trim(unnested) AS genre
FROM
  animes,
  unnest(string_to_array(genre, ',')) AS unnested
ORDER BY
  genre;

-- name: GetRecentlyUpdatedAnimes :many
SELECT
  *
FROM
  animes
WHERE
  animes.mal_id IS NOT NULL
  OR animes.mal_id != 0
ORDER BY
  updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRecentlyUpdatedAnimesCount :one
SELECT
  COUNT(*)
FROM
  animes
WHERE
  animes.mal_id IS NOT NULL
  OR animes.mal_id != 0;

-- name: GetAnimeByGenre :many
SELECT
  *
FROM
  animes
WHERE
  genre ILIKE '%' || sqlc.arg(genre) || '%'
ORDER BY
  updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAnimeByGenreCount :one
SELECT
  COUNT(*)
FROM
  animes
WHERE
  genre ILIKE '%' || sqlc.arg(genre) || '%';

-- name: GetRandomAnime :one
SELECT
  *
FROM
  animes
WHERE
  animes.mal_id IS NOT NULL
  OR animes.mal_id != 0
ORDER BY
  RANDOM()
LIMIT 1;

-- name: GetRandomAnimeByGenre :one
SELECT
  *
FROM
  animes
WHERE
  genre ILIKE '%' || sqlc.arg(genre) || '%'
  AND (animes.mal_id IS NOT NULL
    OR animes.mal_id != 0)
ORDER BY
  RANDOM()
LIMIT 1;

-- name: GetAnimeById :one
SELECT
  *
FROM
  animes
WHERE
  id = $1;

-- name: GetAnimeByMalId :one
SELECT
  *
FROM
  animes
WHERE
  mal_id = $1;

-- name: GetAnimeByHiAnimeId :one
SELECT
  *
FROM
  animes
WHERE
  hi_anime_id = $1;

-- name: GetAnimeMetadataByMalId :one
SELECT
  *
FROM
  anime_metadata
WHERE
  mal_id = $1;

-- name: GetAnimesByMalIds :many
SELECT
  *
FROM
  animes
WHERE
  mal_id = ANY (sqlc.arg(mal_ids)::int[])
ORDER BY
  updated_at DESC;

-- name: GetAnimesByHiAnimeIds :many
SELECT
  *
FROM
  animes
WHERE
  hi_anime_id = ANY (sqlc.arg(hi_anime_ids)::text[])
ORDER BY
  updated_at DESC;

-- name: SearchAnimes :many
SELECT
  animes.*,
  ts_rank(animes.search_vector, plainto_tsquery(sqlc.arg(query))) AS query_rank
FROM
  animes
WHERE (sqlc.arg(query) = ''
  OR sqlc.arg(query) IS NULL
  OR ename % sqlc.arg(query)
  OR jname % sqlc.arg(query)
  OR search_vector @@ plainto_tsquery('english', sqlc.arg(query)))
AND (sqlc.arg(genre) = ''
  OR sqlc.arg(genre) IS NULL
  OR genre ILIKE '%' || sqlc.arg(genre) || '%')
AND animes.mal_id IS NOT NULL
ORDER BY
  query_rank DESC
LIMIT $1 OFFSET $2;

-- name: SearchAnimesCount :one
SELECT
  COUNT(*)
FROM
  animes
WHERE (sqlc.arg(query) = ''
  OR sqlc.arg(query) IS NULL
  OR ename % sqlc.arg(query)
  OR jname % sqlc.arg(query)
  OR search_vector @@ plainto_tsquery('english', sqlc.arg(query)))
AND (sqlc.arg(genre) = ''
  OR sqlc.arg(genre) IS NULL
  OR genre ILIKE '%' || sqlc.arg(genre) || '%')
AND animes.mal_id IS NOT NULL;

-- name: InsertAnime :exec
INSERT INTO animes(ename, jname, image_url, genre, hi_anime_id, mal_id, anilist_id, last_episode, created_at, updated_at, season, season_year)
  VALUES (sqlc.arg(ename), sqlc.arg(jname), sqlc.arg(image_url), sqlc.arg(genre), sqlc.arg(hi_anime_id), sqlc.arg(mal_id), sqlc.arg(anilist_id), sqlc.arg(last_episode), COALESCE(sqlc.arg(created_at), NOW()), COALESCE(sqlc.arg(updated_at), NOW()), sqlc.arg(season), sqlc.arg(season_year))
RETURNING
  *;

-- name: InsertMultipleAnimes :copyfrom
INSERT INTO animes(ename, jname, image_url, genre, hi_anime_id, mal_id, anilist_id, last_episode, season, season_year)
  VALUES (sqlc.arg(ename), sqlc.arg(jname), sqlc.arg(image_url), sqlc.arg(genre), sqlc.arg(hi_anime_id), sqlc.arg(mal_id), sqlc.arg(anilist_id), sqlc.arg(last_episode), sqlc.arg(season), sqlc.arg(season_year));

-- name: UpdateAnime :exec
UPDATE
  animes
SET
  ename = sqlc.arg(ename),
  jname = sqlc.arg(jname),
  image_url = sqlc.arg(image_url),
  genre = sqlc.arg(genre),
  hi_anime_id = sqlc.arg(hi_anime_id),
  mal_id = sqlc.arg(mal_id),
  anilist_id = sqlc.arg(anilist_id),
  last_episode = sqlc.arg(last_episode),
  updated_at = COALESCE(sqlc.arg(updated_at), NOW()),
  season = sqlc.arg(season),
  season_year = sqlc.arg(season_year)
WHERE
  id = sqlc.arg(id)
RETURNING
  *;

-- name: UpsertAnimeMetadata :exec
INSERT INTO anime_metadata(mal_id, description, main_picture_url, media_type, rating, airing_status, avg_episode_duration, total_episodes, studio, rank, mean, scoringUsers, popularity, airing_start_date, airing_end_date, source, trailer_embed_url, season_year, season)
  VALUES (sqlc.arg(mal_id), sqlc.arg(description), sqlc.arg(main_picture_url), sqlc.arg(media_type), sqlc.arg(rating), sqlc.arg(airing_status), sqlc.arg(avg_episode_duration), sqlc.arg(total_episodes), sqlc.arg(studio), sqlc.arg(rank), sqlc.arg(mean), sqlc.arg(scoringUsers), sqlc.arg(popularity), sqlc.arg(airing_start_date), sqlc.arg(airing_end_date), sqlc.arg(source), sqlc.arg(trailer_embed_url), sqlc.arg(season_year), sqlc.arg(season))
ON CONFLICT (mal_id)
  DO UPDATE SET
    description = EXCLUDED.description,
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
  RETURNING
    *;

-- name: UpdateAnimeMetadataTrailer :exec
UPDATE
  anime_metadata
SET
  trailer_embed_url = sqlc.arg(trailer_embed_url),
  updated_at = NOW()
WHERE
  mal_id = sqlc.arg(mal_id)
RETURNING
  *;

-- name: GetAnimeBySeasonAndYear :many
SELECT
  *
FROM
  animes
WHERE
  season = @season::season
  AND season_year = @season_year::int
ORDER BY
  updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAnimeBySeasonAndYearCount :one
SELECT
  COUNT(*)
FROM
  animes
WHERE
  season = @season::season
  AND season_year = @season_year::int;

-- name: GetAnimeByYear :many
SELECT
  *
FROM
  animes
WHERE
  season_year = @season_year::int
ORDER BY
  updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAnimeByYearCount :one
SELECT
  COUNT(*)
FROM
  animes
WHERE
  season_year = @season_year::int;

-- name: GetAnimeBySeason :many
SELECT
  *
FROM
  animes
WHERE
  season = @season::season
ORDER BY
  season_year DESC
LIMIT $1 OFFSET $2;

-- name: GetAnimeBySeasonCount :one
SELECT
  COUNT(*)
FROM
  animes
WHERE
  season = @season::season;

-- name: GetAnimeCatalog :many
WITH p AS (
  SELECT
    -- normalized/trimmed search
    NULLIF(trim(sqlc.narg(search)::text), '') AS q,
    -- normalized genres (lowercased, trimmed) or NULL
    CASE WHEN sqlc.narg(genres)::text[] IS NULL THEN
      NULL
    ELSE
      (
        SELECT
          array_agg(lower(trim(g)))
        FROM
          unnest(sqlc.narg(genres)::text[]) AS u(g)
        WHERE
          trim(g) <> '')
    END AS g,
    sqlc.narg(genres_mode)::text AS gm,
    sqlc.narg(sort_by)::text AS sb,
    sqlc.narg(sort_order)::text AS so
)
SELECT
  a.*,
  l.id AS library_id,
  l.user_id AS library_user_id,
  l.anime_id AS library_anime_id,
  l.status AS library_status,
  l.watched_episodes AS library_watched_episodes,
  l.created_at AS library_created_at,
  l.updated_at AS library_updated_at,
  CASE WHEN p.q IS NOT NULL THEN
    ts_rank(a.search_vector, plainto_tsquery('english', p.q))
  ELSE
    NULL
  END AS query_rank
FROM
  animes a
  -- Only JOIN library when user_id is provided
  LEFT JOIN library l ON (sqlc.narg(user_id)::varchar IS NOT NULL
      AND a.id = l.anime_id
      AND l.user_id = sqlc.narg(user_id)::varchar)
  CROSS JOIN p
WHERE
  -- only MAL-linked rows
(a.mal_id IS NOT NULL
    AND a.mal_id <> 0)
  -- search (skip when q is null)
  AND (p.q IS NULL
    OR a.ename % p.q
    OR a.jname % p.q
    OR a.search_vector @@ plainto_tsquery('english', p.q))
  -- seasons (skip when null)
  AND (sqlc.narg(seasons)::text[] IS NULL
    OR a.season = ANY (sqlc.narg(seasons)::season[]))
  -- years list (skip when null)
  AND (sqlc.narg(years)::int[] IS NULL
    OR a.season_year = ANY (sqlc.narg(years)::int[]))
  -- year range (skip each bound when null)
  AND (sqlc.narg(year_min)::int IS NULL
    OR a.season_year >= sqlc.narg(year_min)::int)
  AND (sqlc.narg(year_max)::int IS NULL
    OR a.season_year <= sqlc.narg(year_max)::int)
  -- genres ANY/ALL using generated genres_arr (skip when null/empty)
  AND (p.g IS NULL
    OR (
      CASE WHEN p.gm = 'all' THEN
        a.genres_arr @> p.g
      ELSE
        a.genres_arr && p.g
      END))
    -- Library-only filtering (when user_id provided, only show library entries)
    AND (sqlc.narg(user_id)::varchar IS NULL -- catalog mode
      OR l.user_id IS NOT NULL) -- library mode (must be in library)
    -- Library status filtering
    AND (sqlc.narg(library_status)::library_status IS NULL
      OR l.status = sqlc.narg(library_status)::library_status)
  ORDER BY
    -- relevance
    CASE WHEN p.sb = 'relevance'
      AND p.so = 'asc' THEN
      CASE WHEN p.q IS NOT NULL THEN
        ts_rank(a.search_vector, plainto_tsquery('english', p.q))
      END
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'relevance'
      AND p.so = 'desc' THEN
      CASE WHEN p.q IS NOT NULL THEN
        ts_rank(a.search_vector, plainto_tsquery('english', p.q))
      END
    END DESC NULLS LAST,
    -- ename
    CASE WHEN p.sb = 'ename'
      AND p.so = 'asc' THEN
      a.ename
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'ename'
      AND p.so = 'desc' THEN
      a.ename
    END DESC NULLS LAST,
    -- jname
    CASE WHEN p.sb = 'jname'
      AND p.so = 'asc' THEN
      a.jname
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'jname'
      AND p.so = 'desc' THEN
      a.jname
    END DESC NULLS LAST,
    -- season
    CASE WHEN p.sb = 'season'
      AND p.so = 'asc' THEN
      a.season::text
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'season'
      AND p.so = 'desc' THEN
      a.season::text
    END DESC NULLS LAST,
    -- year
    CASE WHEN p.sb = 'year'
      AND p.so = 'asc' THEN
      a.season_year
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'year'
      AND p.so = 'desc' THEN
      a.season_year
    END DESC NULLS LAST,
    -- anime updated_at
    CASE WHEN p.sb = 'anime_updated_at'
      AND p.so = 'asc' THEN
      a.updated_at
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'anime_updated_at'
      AND p.so = 'desc' THEN
      a.updated_at
    END DESC NULLS LAST,
    -- library updated_at (only when library is joined)
    CASE WHEN p.sb = 'library_updated_at'
      AND p.so = 'asc' THEN
      l.updated_at
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'library_updated_at'
      AND p.so = 'desc' THEN
      l.updated_at
    END DESC NULLS LAST,
    -- legacy updated_at (maps to anime_updated_at for backward compatibility)
    CASE WHEN p.sb = 'updated_at'
      AND p.so = 'asc' THEN
      a.updated_at
    END ASC NULLS LAST,
    CASE WHEN p.sb = 'updated_at'
      AND p.so = 'desc' THEN
      a.updated_at
    END DESC NULLS LAST,
    -- stable tiebreakers (helpful for deterministic paging)
    a.updated_at DESC,
    a.id DESC
  LIMIT $1 OFFSET $2;

-- name: GetAnimeCatalogCount :one
WITH p AS (
  SELECT
    NULLIF(trim(sqlc.narg(search)::text), '') AS q,
    CASE WHEN sqlc.narg(genres)::text[] IS NULL THEN
      NULL
    ELSE
      (
        SELECT
          array_agg(lower(trim(g)))
        FROM
          unnest(sqlc.narg(genres)::text[]) AS u(g)
        WHERE
          trim(g) <> '')
    END AS g,
    sqlc.narg(genres_mode)::text AS gm
)
SELECT
  COUNT(*)
FROM
  animes a
  -- Only JOIN library when user_id is provided
  LEFT JOIN library l ON (sqlc.narg(user_id)::varchar IS NOT NULL
      AND a.id = l.anime_id
      AND l.user_id = sqlc.narg(user_id)::varchar)
  CROSS JOIN p
WHERE
  -- only MAL-linked rows
(a.mal_id IS NOT NULL
    AND a.mal_id <> 0)
  -- search (skip when q is null)
  AND (p.q IS NULL
    OR a.ename % p.q
    OR a.jname % p.q
    OR a.search_vector @@ plainto_tsquery('english', p.q))
  -- seasons (skip when null)
  AND (sqlc.narg(seasons)::text[] IS NULL
    OR a.season = ANY (sqlc.narg(seasons)::season[]))
  -- years list (skip when null)
  AND (sqlc.narg(years)::int[] IS NULL
    OR a.season_year = ANY (sqlc.narg(years)::int[]))
  -- year range (skip each bound when null)
  AND (sqlc.narg(year_min)::int IS NULL
    OR a.season_year >= sqlc.narg(year_min)::int)
  AND (sqlc.narg(year_max)::int IS NULL
    OR a.season_year <= sqlc.narg(year_max)::int)
  -- genres ANY/ALL using generated genres_arr (skip when null/empty)
  AND (p.g IS NULL
    OR (
      CASE WHEN p.gm = 'all' THEN
        a.genres_arr @> p.g
      ELSE
        a.genres_arr && p.g
      END))
    -- Library-only filtering (when user_id provided, only show library entries)
    AND (sqlc.narg(user_id)::varchar IS NULL -- catalog mode
      OR l.user_id IS NOT NULL) -- library mode (must be in library)
    -- Library status filtering
    AND (sqlc.narg(library_status)::library_status IS NULL
      OR l.status = sqlc.narg(library_status)::library_status);

-- name: GetGenrePreviews :many
WITH g AS (
  SELECT DISTINCT
    unnest(a.genres_arr) AS genre
  FROM
    animes a
)
SELECT
  g.genre::text AS name,
  COALESCE(ARRAY (
      SELECT
        a2.image_url::text
      FROM animes a2
      WHERE
        a2.genres_arr @> ARRAY[g.genre]::text[] ORDER BY a2.season_year DESC, a2.updated_at DESC, a2.id DESC LIMIT 6), ARRAY[]::text[]) AS previews
FROM
  g
ORDER BY
  g.genre;

