-- name: ListActiveWorkTypes :many
SELECT * FROM work_types WHERE active = TRUE ORDER BY sort_order;

-- name: GetWorkTypeByID :one
SELECT * FROM work_types WHERE id = $1;

-- name: GetWorkTypeBySlug :one
SELECT * FROM work_types WHERE slug = $1;
