-- name: CreateResetPasswordToken :one
INSERT INTO reset_password_tokens(user_id)
  VALUES (sqlc.arg(user_id))
RETURNING
  *;

-- name: DeleteResetPasswordToken :exec
DELETE FROM reset_password_tokens
WHERE token = sqlc.arg(token);

-- name: GetResetPasswordToken :one
SELECT
  sqlc.embed(users),
  sqlc.embed(reset_password_tokens)
FROM
  reset_password_tokens
  INNER JOIN users ON users.id = reset_password_tokens.user_id
WHERE
  reset_password_tokens.token = sqlc.arg(token)
  AND reset_password_tokens.expires_at > NOW();

