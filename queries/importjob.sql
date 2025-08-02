-- name: CreateLibraryImportJob :one
INSERT INTO library_import_jobs(user_id, provider)
  VALUES (sqlc.arg(user_id), sqlc.arg(provider))
RETURNING
  id;

-- name: GetLibraryImportJob :one
SELECT
  *
FROM
  library_import_jobs
WHERE
  id = sqlc.arg(id);

-- name: GetLibraryImportJobByUserId :many
SELECT
  *
FROM
  library_import_jobs
WHERE
  user_id = sqlc.arg(user_id)
ORDER BY
  created_at DESC;

-- name: UpdateLibraryImportJob :exec
UPDATE
  library_import_jobs
SET
  status = sqlc.arg(status)::library_import_status,
  error_message = sqlc.arg(error_message),
  updated_at = NOW(),
  completed_at = CASE WHEN sqlc.arg(status)::library_import_status IN ('completed', 'failed') THEN
    NOW()
  ELSE
    NULL
  END
WHERE
  id = sqlc.arg(id);

