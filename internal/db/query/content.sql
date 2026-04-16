-- name: GetContentByID :one
SELECT * FROM content WHERE id = $1 AND deleted_at IS NULL;

-- name: ListContent :many
SELECT * FROM content WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListContentByStatus :many
SELECT * FROM content WHERE status = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListContentByCompany :many
SELECT * FROM content WHERE company_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountContent :one
SELECT COUNT(*) FROM content WHERE deleted_at IS NULL;

-- name: CountContentByStatus :many
SELECT status, COUNT(*) AS count FROM content WHERE deleted_at IS NULL GROUP BY status;

-- name: CreateContent :one
INSERT INTO content (title, content_type, platforms, topics, company_id, author_id, campaign_id, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateContent :one
UPDATE content SET
    title = $2, content_type = $3, platforms = $4, topics = $5,
    company_id = $6, status = $7, notes = $8
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: ReviewContent :one
UPDATE content SET status = $2, review_notes = $3
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: PublishContent :one
UPDATE content SET status = 'published', publish_date = $2, published_url = $3
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteContent :exec
UPDATE content SET deleted_at = NOW() WHERE id = $1;
