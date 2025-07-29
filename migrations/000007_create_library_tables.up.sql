CREATE TYPE library_status AS ENUM(
  'planning',
  'watching',
  'completed',
  'dropped',
  'paused'
);

CREATE TABLE library (
  id VARCHAR(21) PRIMARY KEY DEFAULT generate_nanoid (),
  user_id VARCHAR(21) NOT NULL,
  anime_id VARCHAR(21) NOT NULL,
  status library_status NOT NULL DEFAULT 'planning',
  watched_episodes INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, anime_id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (anime_id) REFERENCES animes (id) ON DELETE CASCADE
);

CREATE TYPE library_actions AS ENUM(
  'add_entry',
  'update_progress',
  'update_status',
  'delete_entry'
);

CREATE TYPE library_sync_status AS ENUM('pending', 'success', 'failed', 'skipped');

CREATE TABLE external_library_sync (
  user_id VARCHAR(21) NOT NULL,
  anime_id VARCHAR(21) NOT NULL,
  provider Provider NOT NULL,
  action library_actions NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}',
  status library_sync_status NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, anime_id, provider, action),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (anime_id) REFERENCES animes (id) ON DELETE CASCADE
);
