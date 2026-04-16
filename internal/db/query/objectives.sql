-- name: GetObjectiveByID :one
SELECT * FROM objectives WHERE id = $1 AND deleted_at IS NULL;

-- name: ListObjectivesByCompany :many
SELECT * FROM objectives WHERE company_id = $1 AND deleted_at IS NULL ORDER BY quarter DESC;

-- name: ListObjectivesByQuarter :many
SELECT * FROM objectives WHERE quarter = $1 AND deleted_at IS NULL ORDER BY company_id;

-- name: CreateObjective :one
INSERT INTO objectives (title, description, company_id, owner_id, quarter, year, key_results, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateObjective :one
UPDATE objectives SET
    title = $2, description = $3, status = $4, progress = $5, key_results = $6, notes = $7
WHERE id = $1 AND deleted_at IS NULL RETURNING *;
