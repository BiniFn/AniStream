-- name: CreateSession :one
INSERT INTO
  sessions (user_id)
VALUES
  (sqlc.arg (user_id))
RETURNING
  *;

-- name: GetUserBySessionID :one
SELECT
  users.*
FROM
  sessions
  INNER JOIN users ON sessions.user_id = users.id
WHERE
  sessions.id = sqlc.arg (id);

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
  id = sqlc.arg (id);
