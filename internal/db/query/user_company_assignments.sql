-- name: ListAssignmentsByCompany :many
SELECT uca.*, u.full_name, u.email, u.role AS user_role, u.status AS user_status
FROM user_company_assignments uca
JOIN users u ON u.id = uca.user_id
WHERE uca.company_id = $1 AND (uca.end_date IS NULL OR uca.end_date > CURRENT_DATE)
ORDER BY u.full_name;

-- name: ListAssignmentsByUser :many
SELECT uca.*, c.name AS company_name, c.short_code
FROM user_company_assignments uca
JOIN companies c ON c.id = uca.company_id
WHERE uca.user_id = $1 AND (uca.end_date IS NULL OR uca.end_date > CURRENT_DATE)
ORDER BY c.name;

-- name: CreateAssignment :one
INSERT INTO user_company_assignments (user_id, company_id, role_in_company, can_view, can_edit, can_approve)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: DeleteAssignment :exec
DELETE FROM user_company_assignments WHERE id = $1;
