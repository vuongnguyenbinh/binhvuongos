-- name: GetInboxItemByID :one
SELECT * FROM inbox_items WHERE id = $1;

-- name: ListInboxItems :many
SELECT * FROM inbox_items ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListInboxItemsByStatus :many
SELECT * FROM inbox_items WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountInboxItems :one
SELECT COUNT(*) FROM inbox_items;

-- name: CountInboxItemsByStatus :one
SELECT COUNT(*) FROM inbox_items WHERE status = $1;

-- name: CreateInboxItem :one
INSERT INTO inbox_items (content, url, source, item_type, company_id, submitted_by, attachments, external_ref)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetInboxByExternalRef :one
SELECT * FROM inbox_items
WHERE source = $1 AND external_ref = $2
LIMIT 1;

-- name: UpdateInboxItem :one
UPDATE inbox_items SET content = $2, url = $3, item_type = $4, company_id = $5
WHERE id = $1 RETURNING *;

-- name: TriageInboxItem :one
UPDATE inbox_items SET
    status = 'done', destination = $2, triage_notes = $3,
    converted_to_type = $4, converted_to_id = $5, processed_at = NOW()
WHERE id = $1 RETURNING *;

-- name: ArchiveOldInboxItems :exec
UPDATE inbox_items SET status = 'archived'
WHERE status = 'raw' AND created_at < NOW() - INTERVAL '7 days';
