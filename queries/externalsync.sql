-- name: UpsertLibrarySync :exec
INSERT INTO external_library_sync(user_id, anime_id, provider, action, payload)
  VALUES (sqlc.arg(user_id), sqlc.arg(anime_id), sqlc.arg(provider), sqlc.arg(action), sqlc.arg(payload))
ON CONFLICT (user_id, anime_id, provider, action)
  DO UPDATE SET
    payload = EXCLUDED.payload,
    status = 'pending',
    updated_at = NOW();

-- name: UpdateLibrarySyncStatus :exec
UPDATE
  external_library_sync
SET
  status = sqlc.arg(status),
  updated_at = NOW()
WHERE
  user_id = sqlc.arg(user_id)
  AND anime_id = sqlc.arg(anime_id)
  AND provider = sqlc.arg(provider)
  AND action = sqlc.arg(action);

-- name: GetPendingLibrarySyncs :many
SELECT
  *
FROM
  external_library_sync
WHERE
  user_id = sqlc.arg(user_id)
  AND status = 'pending';

-- name: GetFailedLibrarySyncs :many
SELECT
  *
FROM
  external_library_sync
WHERE
  user_id = sqlc.arg(user_id)
  AND status = 'failed';

-- name: GetAllLibrarySyncsForAnime :many
SELECT
  *
FROM
  external_library_sync
WHERE
  user_id = sqlc.arg(user_id)
  AND anime_id = sqlc.arg(anime_id);

-- name: DeleteLibrarySync :exec
DELETE FROM external_library_sync
WHERE user_id = sqlc.arg(user_id)
  AND anime_id = sqlc.arg(anime_id)
  AND provider = sqlc.arg(provider)
  AND action = sqlc.arg(action);

-- name: GetPendingLibrarySyncsForAllUsers :many
SELECT
  *
FROM
  external_library_sync
WHERE
  status = 'pending'
ORDER BY
  updated_at ASC;

