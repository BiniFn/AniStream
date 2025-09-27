CREATE OR REPLACE FUNCTION notify_library_sync()
  RETURNS TRIGGER
  AS $$
BEGIN
  PERFORM
    pg_notify('library_sync', json_build_object('user_id', NEW.user_id, 'anime_id', NEW.anime_id, 'provider', NEW.provider, 'action', NEW.action, 'payload', NEW.payload)::text);
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER sync_insert_update_notify
  AFTER INSERT OR UPDATE ON external_library_sync
  FOR EACH ROW
  WHEN(NEW.status = 'pending')
  EXECUTE FUNCTION notify_library_sync();

