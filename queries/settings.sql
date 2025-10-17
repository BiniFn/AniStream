-- name: GetSettingsOfUser :one
SELECT
  sqlc.embed(settings),
  sqlc.embed(themes)
FROM
  settings
  INNER JOIN themes ON settings.theme_id = themes.id
WHERE
  settings.user_id = sqlc.arg(user_id);

-- name: SaveSettings :one
WITH upserted AS (
INSERT INTO settings(user_id, auto_next_episode, auto_play_episode, auto_resume_episode, incognito_mode, theme_id)
    VALUES (sqlc.arg(user_id), sqlc.arg(auto_next_episode), sqlc.arg(auto_play_episode), sqlc.arg(auto_resume_episode), sqlc.arg(incognito_mode), sqlc.arg(theme_id))
  ON CONFLICT (user_id)
    DO UPDATE SET
      auto_next_episode = EXCLUDED.auto_next_episode,
      auto_play_episode = EXCLUDED.auto_play_episode,
      auto_resume_episode = EXCLUDED.auto_resume_episode,
      incognito_mode = EXCLUDED.incognito_mode,
      theme_id = EXCLUDED.theme_id
    RETURNING
      *
)
  SELECT
    upserted.*,
    sqlc.embed(themes)
  FROM
    upserted
  LEFT JOIN themes ON themes.id = upserted.theme_id;

-- name: ListThemes :many
SELECT
  *
FROM
  themes
ORDER BY
  id ASC;

