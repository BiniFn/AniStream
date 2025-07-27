-- name: CreateResetPasswordToken :one
INSERT INTO
  reset_password_tokens (user_id)
VALUES
  (sqlc.arg (user_id))
RETURNING
  *;

-- name: GetResetPasswordToken :one
SELECT
  *
FROM
  reset_password_tokens
WHERE
  token = sqlc.arg (token);

-- name: DeleteResetPasswordToken :exec
DELETE FROM reset_password_tokens
WHERE
  token = sqlc.arg (token);

-- name: GetUserByResetPasswordToken :one
SELECT
  users.*
FROM
  reset_password_tokens
  INNER JOIN users on users.id = reset_password_tokens.user_id
WHERE
  reset_password_tokens.token = sqlc.arg (token);
