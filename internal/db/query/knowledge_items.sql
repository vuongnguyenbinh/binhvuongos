-- name: GetKnowledgeItemByID :one
SELECT * FROM knowledge_items WHERE id = $1 AND deleted_at IS NULL;

-- name: ListKnowledgeItems :many
SELECT * FROM knowledge_items WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListKnowledgeItemsByCategory :many
SELECT * FROM knowledge_items WHERE category = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountKnowledgeItems :one
SELECT COUNT(*) FROM knowledge_items WHERE deleted_at IS NULL;

-- name: SearchKnowledgeItems :many
SELECT * FROM knowledge_items
WHERE deleted_at IS NULL
  AND to_tsvector('simple', immutable_unaccent(title || ' ' || COALESCE(description, ''))) @@ plainto_tsquery('simple', immutable_unaccent($1))
ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CreateKnowledgeItem :one
INSERT INTO knowledge_items (title, description, body, category, topics, quality_rating, scope, company_id, format, source_url, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: UpdateKnowledgeItem :one
UPDATE knowledge_items SET
    title = $2, description = $3, body = $4, category = $5,
    topics = $6, quality_rating = $7, scope = $8, format = $9, source_url = $10
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteKnowledgeItem :exec
UPDATE knowledge_items SET deleted_at = NOW() WHERE id = $1;
