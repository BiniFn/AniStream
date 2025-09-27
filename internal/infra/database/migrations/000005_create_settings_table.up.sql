CREATE TABLE settings(
  user_id varchar(21) PRIMARY KEY,
  auto_next_episode boolean NOT NULL DEFAULT TRUE,
  auto_play_episode boolean NOT NULL DEFAULT TRUE,
  auto_resume_episode boolean NOT NULL DEFAULT FALSE,
  incognito_mode boolean NOT NULL DEFAULT FALSE)
