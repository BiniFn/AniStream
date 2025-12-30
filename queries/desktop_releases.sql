-- name: GetAllDesktopReleases :many
SELECT
    *
FROM
    desktop_releases
ORDER BY
    version DESC,
    platform ASC;

-- name: GetDesktopReleasesByVersion :many
SELECT
    *
FROM
    desktop_releases
WHERE
    version = sqlc.arg (version)
ORDER BY
    platform ASC;

-- name: GetLatestDesktopReleases :many
SELECT
    *
FROM
    desktop_releases
WHERE
    version = (
        SELECT
            version
        FROM
            desktop_releases
        ORDER BY
            created_at DESC
        LIMIT 1)
ORDER BY
    platform ASC;

-- name: GetDesktopReleaseByVersionAndPlatform :one
SELECT
    *
FROM
    desktop_releases
WHERE
    version = sqlc.arg (version)
    AND platform = sqlc.arg (platform);

-- name: InsertDesktopRelease :one
INSERT INTO desktop_releases (version, platform, download_url, file_name, file_size, release_notes)
    VALUES (sqlc.arg (version), sqlc.arg (platform), sqlc.arg (download_url), sqlc.arg (file_name), sqlc.arg (file_size), sqlc.arg (release_notes))
RETURNING
    *;

-- name: DeleteDesktopReleasesByVersion :exec
DELETE FROM desktop_releases
WHERE version = sqlc.arg (version);

