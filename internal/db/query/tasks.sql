-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTasks :many
SELECT * FROM tasks WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListTasksByStatus :many
SELECT * FROM tasks WHERE status = $1 AND deleted_at IS NULL ORDER BY priority DESC, due_date ASC NULLS LAST;

-- name: ListTasksByCompany :many
SELECT * FROM tasks WHERE company_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListTasksByAssignee :many
SELECT * FROM tasks WHERE assignee_id = $1 AND deleted_at IS NULL ORDER BY due_date ASC NULLS LAST;

-- name: ListTasksDueToday :many
SELECT * FROM tasks
WHERE due_date <= CURRENT_DATE AND status NOT IN ('done', 'cancelled') AND deleted_at IS NULL
ORDER BY priority DESC, due_date ASC;

-- name: CountTasksByStatus :many
SELECT status, COUNT(*) AS count FROM tasks WHERE deleted_at IS NULL GROUP BY status;

-- name: CountOverdueTasks :one
SELECT COUNT(*) FROM tasks
WHERE due_date < CURRENT_DATE AND status NOT IN ('done', 'cancelled') AND deleted_at IS NULL;

-- name: CreateTask :one
INSERT INTO tasks (title, description, category, group_name, company_id, assignee_id, campaign_id, status, priority, due_date, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: UpdateTask :one
UPDATE tasks SET
    title = $2, description = $3, category = $4, group_name = $5,
    company_id = $6, assignee_id = $7, status = $8, priority = $9, due_date = $10
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: UpdateTaskStatus :one
UPDATE tasks SET status = $2,
    started_at = CASE WHEN $2 = 'in_progress' AND started_at IS NULL THEN NOW() ELSE started_at END,
    completed_at = CASE WHEN $2 = 'done' THEN NOW() ELSE completed_at END
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteTask :exec
UPDATE tasks SET deleted_at = NOW() WHERE id = $1;
