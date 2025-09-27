ALTER TABLE sessions
  ALTER COLUMN expires_at SET DEFAULT (CURRENT_TIMESTAMP + interval '1 month');

