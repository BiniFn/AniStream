-- name: GetSettingsOfUser :one
SELECT
  *
FROM
  settings
WHERE
  user_id = sqlc.arg (user_id);

-- name: SaveSettings :one
INSERT INTO
  settings (
    user_id,
    auto_next_episode,
    auto_play_episode,
    auto_resume_episode,
    incognito_mode
  )
VALUES
  (
    sqlc.arg (user_id),
    sqlc.arg (auto_next_episode),
    sqlc.arg (auto_play_episode),
    sqlc.arg (auto_resume_episode),
    sqlc.arg (incognito_mode)
  )
ON CONFLICT (user_id) DO UPDATE
SET
  auto_next_episode = EXCLUDED.auto_next_episode,
  auto_play_episode = EXCLUDED.auto_play_episode,
  auto_resume_episode = EXCLUDED.auto_resume_episode,
  incognito_mode = EXCLUDED.incognito_mode
RETURNING
  *;
