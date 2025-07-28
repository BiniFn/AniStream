CREATE TABLE settings (
  user_id VARCHAR(21) PRIMARY KEY,
  auto_next_episode BOOLEAN NOT NULL DEFAULT true,
  auto_play_episode BOOLEAN NOT NULL DEFAULT true,
  auto_resume_episode BOOLEAN NOT NULL DEFAULT false,
  incognito_mode BOOLEAN NOT NULL DEFAULT false
)
