-- name: GetBookmarkByID :one
SELECT * FROM bookmarks WHERE id = $1 AND deleted_at IS NULL;

-- name: ListBookmarks :many
SELECT * FROM bookmarks WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountBookmarks :one
SELECT COUNT(*) FROM bookmarks WHERE deleted_at IS NULL;

-- name: CreateBookmark :one
INSERT INTO bookmarks (title, url, description, tags, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateBookmark :one
UPDATE bookmarks SET title = $2, url = $3, description = $4, tags = $5, notes = $6
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteBookmark :exec
UPDATE bookmarks SET deleted_at = NOW() WHERE id = $1;
