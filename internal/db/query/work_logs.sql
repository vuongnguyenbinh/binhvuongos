-- name: GetWorkLogByID :one
SELECT * FROM work_logs WHERE id = $1;

-- name: ListWorkLogs :many
SELECT * FROM work_logs ORDER BY work_date DESC, created_at DESC LIMIT $1 OFFSET $2;

-- name: ListWorkLogsByStatus :many
SELECT * FROM work_logs WHERE status = $1 ORDER BY work_date DESC LIMIT $2 OFFSET $3;

-- name: ListWorkLogsByUser :many
SELECT * FROM work_logs WHERE user_id = $1 ORDER BY work_date DESC LIMIT $2 OFFSET $3;

-- name: ListWorkLogsByCompany :many
SELECT * FROM work_logs WHERE company_id = $1 ORDER BY work_date DESC LIMIT $2 OFFSET $3;

-- name: CountWorkLogs :one
SELECT COUNT(*) FROM work_logs;

-- name: CountPendingWorkLogs :one
SELECT COUNT(*) FROM work_logs WHERE status = 'submitted';

-- name: CreateWorkLog :one
INSERT INTO work_logs (work_date, user_id, company_id, work_type_id, campaign_id, quantity, sheet_url, evidence_url, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateWorkLog :one
UPDATE work_logs SET
    work_date = $2, work_type_id = $3, quantity = $4,
    sheet_url = $5, evidence_url = $6, notes = $7
WHERE id = $1 RETURNING *;

-- name: ApproveWorkLog :one
UPDATE work_logs SET status = 'approved', reviewed_at = NOW(), reviewed_by = $2, admin_notes = $3
WHERE id = $1 RETURNING *;

-- name: RejectWorkLog :one
UPDATE work_logs SET status = 'rejected', reviewed_at = NOW(), reviewed_by = $2, admin_notes = $3
WHERE id = $1 RETURNING *;

-- name: BatchApproveWorkLogs :exec
UPDATE work_logs SET status = 'approved', reviewed_at = NOW(), reviewed_by = $2
WHERE id = ANY($1::UUID[]);
