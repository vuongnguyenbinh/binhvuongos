-- name: GetCampaignByID :one
SELECT * FROM campaigns WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCampaigns :many
SELECT * FROM campaigns WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListCampaignsByCompany :many
SELECT * FROM campaigns WHERE company_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: ListCampaignsByStatus :many
SELECT * FROM campaigns WHERE status = $1 AND deleted_at IS NULL ORDER BY start_date DESC;

-- name: CountCampaigns :one
SELECT COUNT(*) FROM campaigns WHERE deleted_at IS NULL;

-- name: CreateCampaign :one
INSERT INTO campaigns (name, description, company_id, owner_id, campaign_type, status, start_date, end_date, target_json, budget, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: UpdateCampaign :one
UPDATE campaigns SET
    name = $2, description = $3, campaign_type = $4, status = $5,
    start_date = $6, end_date = $7, target_json = $8, budget = $9
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteCampaign :exec
UPDATE campaigns SET deleted_at = NOW() WHERE id = $1;
