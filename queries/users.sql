-- name: GetUserByID :one
SELECT
  *
FROM
  users
WHERE
  id = $1;

-- name: GetUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = $1;

-- name: CreateUser :one
INSERT INTO
  users (username, email, password_hash, profile_picture)
VALUES
  (
    sqlc.arg (username),
    sqlc.arg (email),
    sqlc.arg (password_hash),
    sqlc.arg (profile_picture)
  )
RETURNING
  *;

-- name: UpdateUser :one
UPDATE users
SET
  username = sqlc.arg (username),
  email = sqlc.arg (email),
  profile_picture = sqlc.arg (profile_picture)
WHERE
  id = sqlc.arg (id)
RETURNING
  *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
  id = sqlc.arg (id);

-- name: UpdatePassword :exec
UPDATE users
SET
  password_hash = sqlc.arg (password_hash)
WHERE
  id = sqlc.arg (id);
