CREATE TYPE library_import_status AS ENUM(
  'pending',
  'in_progress',
  'completed',
  'failed'
);

CREATE TABLE library_import_jobs(
  id varchar(21) PRIMARY KEY DEFAULT generate_nanoid(),
  user_id varchar(21) NOT NULL,
  provider Provider NOT NULL,
  status library_import_status NOT NULL DEFAULT 'pending',
  error_message text NULL DEFAULT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at timestamp NULL DEFAULT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_library_import_jobs_user_id ON library_import_jobs(user_id);

CREATE OR REPLACE FUNCTION notify_library_import_job_change()
  RETURNS TRIGGER
  AS $$
DECLARE
  payload json;
BEGIN
  payload = json_build_object('id', NEW.id, 'user_id', NEW.user_id, 'provider', NEW.provider, 'status', NEW.status);
  PERFORM
    pg_notify('library_import_jobs', payload::text);
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER library_import_jobs_notify_trigger
  AFTER INSERT ON library_import_jobs
  FOR EACH ROW
  EXECUTE FUNCTION notify_library_import_job_change();

CREATE OR REPLACE FUNCTION set_updated_at_timestamp()
  RETURNS TRIGGER
  AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER set_library_import_jobs_updated_at
  BEFORE UPDATE ON library_import_jobs
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at_timestamp();

