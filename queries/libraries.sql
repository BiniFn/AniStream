-- name: GetLibrary :many
SELECT
  sqlc.embed(library),
  sqlc.embed(animes)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  library.user_id = sqlc.arg(user_id)
  AND library.status = sqlc.arg(status)
ORDER BY
  library.updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetLibraryCount :one
SELECT
  COUNT(*)
FROM
  library
WHERE
  user_id = sqlc.arg(user_id)
  AND status = sqlc.arg(status);

-- name: IsAnimeInLibrary :one
SELECT
  COUNT(*) > 0
FROM
  library
WHERE
  user_id = sqlc.arg(user_id)
  AND anime_id = sqlc.arg(anime_id);

-- name: GetLibraryByID :one
SELECT
  sqlc.embed(library),
  sqlc.embed(animes)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  library.id = sqlc.arg(id);

-- name: GetLibraryOfUserByAnimeID :one
SELECT
  sqlc.embed(library),
  sqlc.embed(animes)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  library.user_id = sqlc.arg(user_id)
  AND library.anime_id = sqlc.arg(anime_id);

-- name: InsertLibrary :exec
INSERT INTO library(user_id, anime_id, status, watched_episodes, updated_at)
  VALUES (sqlc.arg(user_id), sqlc.arg(anime_id), sqlc.arg(status), sqlc.arg(watched_episodes), coalesce(sqlc.arg(updated_at), NOW()));

-- name: UpdateLibrary :exec
UPDATE
  library
SET
  status = sqlc.arg(status),
  watched_episodes = sqlc.arg(watched_episodes),
  updated_at = coalesce(sqlc.arg(updated_at), NOW())
WHERE
  user_id = sqlc.arg(user_id)
  AND anime_id = sqlc.arg(anime_id);

-- name: DeleteLibrary :exec
DELETE FROM library
WHERE user_id = sqlc.arg(user_id)
  AND anime_id = sqlc.arg(anime_id);

-- name: GetContinueWatchingAnime :many
SELECT
  sqlc.embed(library),
  sqlc.embed(animes)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  status = 'watching'
  AND user_id = sqlc.arg(user_id)
  AND animes.last_episode > library.watched_episodes
ORDER BY
  animes.updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetContinueWatchingAnimeCount :one
SELECT
  COUNT(*)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  status = 'watching'
  AND user_id = sqlc.arg(user_id)
  AND animes.last_episode > library.watched_episodes;

-- name: GetPlanToWatchAnime :many
SELECT
  sqlc.embed(library),
  sqlc.embed(animes)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  status = 'planning'
  AND user_id = sqlc.arg(user_id)
  AND animes.last_episode > library.watched_episodes
ORDER BY
  animes.updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPlanToWatchAnimeCount :one
SELECT
  COUNT(*)
FROM
  library
  INNER JOIN animes ON animes.id = library.anime_id
WHERE
  status = 'planning'
  AND user_id = sqlc.arg(user_id)
  AND animes.last_episode > library.watched_episodes;

