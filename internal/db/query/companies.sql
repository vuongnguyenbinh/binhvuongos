-- name: GetCompanyByID :one
SELECT * FROM companies WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCompanies :many
SELECT * FROM companies WHERE deleted_at IS NULL ORDER BY name LIMIT $1 OFFSET $2;

-- name: ListCompaniesByStatus :many
SELECT * FROM companies WHERE status = $1 AND deleted_at IS NULL ORDER BY name;

-- name: CountCompanies :one
SELECT COUNT(*) FROM companies WHERE deleted_at IS NULL;

-- name: CreateCompany :one
INSERT INTO companies (name, short_code, slug, industry, my_role, scope, status, health, description, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: UpdateCompany :one
UPDATE companies SET
    name = $2, short_code = $3, industry = $4, my_role = $5,
    scope = $6, status = $7, health = $8, description = $9,
    primary_contact_name = $10, primary_contact_phone = $11, primary_contact_email = $12
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: UpdateCompanyHealth :exec
UPDATE companies SET health = $2 WHERE id = $1;

-- name: SoftDeleteCompany :exec
UPDATE companies SET deleted_at = NOW() WHERE id = $1;
