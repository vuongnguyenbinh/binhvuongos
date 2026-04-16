-- name: GetDashboardOutputThisMonth :many
SELECT
    wt.name, wt.unit, wt.icon, wt.slug,
    COALESCE(SUM(wl.quantity), 0)::DECIMAL AS total
FROM work_types wt
LEFT JOIN work_logs wl ON wl.work_type_id = wt.id
    AND wl.work_date >= DATE_TRUNC('month', CURRENT_DATE)
    AND wl.status = 'approved'
WHERE wt.active = TRUE
GROUP BY wt.id, wt.name, wt.unit, wt.icon, wt.slug, wt.sort_order
ORDER BY wt.sort_order;

-- name: GetDashboardCounts :one
SELECT
    (SELECT COUNT(*) FROM work_logs WHERE status = 'submitted') AS pending_reviews,
    (SELECT COUNT(*) FROM content WHERE status = 'review' AND deleted_at IS NULL) AS content_review,
    (SELECT COUNT(*) FROM tasks WHERE due_date < CURRENT_DATE AND status NOT IN ('done', 'cancelled') AND deleted_at IS NULL) AS overdue_tasks,
    (SELECT COUNT(*) FROM inbox_items WHERE status = 'raw') AS raw_inbox,
    (SELECT COUNT(*) FROM tasks WHERE status NOT IN ('done', 'cancelled') AND deleted_at IS NULL) AS open_tasks,
    (SELECT COUNT(*) FROM tasks WHERE status = 'done' AND deleted_at IS NULL) AS done_tasks,
    (SELECT COUNT(*) FROM campaigns WHERE status = 'running' AND deleted_at IS NULL) AS running_campaigns;
