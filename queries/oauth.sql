-- name: GetAllOauthTokensOfUser :many
SELECT
  *
FROM
  oauth_tokens
WHERE
  user_id = sqlc.arg(user_id);

-- name: GetToken :one
SELECT
  *
FROM
  oauth_tokens
WHERE
  user_id = sqlc.arg(user_id)
  AND provider = sqlc.arg(provider);

-- name: GetTokensNearToExpiry :many
SELECT
  *
FROM
  oauth_tokens
WHERE
  expires_at <= NOW() + INTERVAL '10 days'
  AND expires_at > NOW();

-- name: SaveOauthToken :exec
INSERT INTO oauth_tokens(user_id, token, refresh_token, provider, expires_at)
  VALUES (sqlc.arg(user_id), sqlc.arg(token), sqlc.arg(refresh_token), sqlc.arg(provider), sqlc.arg(expires_at));

-- name: UpdateOauthToken :exec
UPDATE
  oauth_tokens
SET
  token = sqlc.arg(token),
  refresh_token = sqlc.arg(refresh_token),
  expires_at = sqlc.arg(expires_at)
WHERE
  user_id = sqlc.arg(user_id)
  AND provider = sqlc.arg(provider);

-- name: DeleteOauthToken :exec
DELETE FROM oauth_tokens
WHERE user_id = sqlc.arg(user_id)
  AND provider = sqlc.arg(provider);

