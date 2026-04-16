package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type DashboardOutput struct {
	Name  string         `json:"name"`
	Unit  string         `json:"unit"`
	Icon  string         `json:"icon"`
	Slug  string         `json:"slug"`
	Total pgtype.Numeric `json:"total"`
}

func (q *Queries) GetDashboardOutputThisMonth(ctx context.Context) ([]DashboardOutput, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT wt.name, wt.unit, COALESCE(wt.icon, '') AS icon, wt.slug,
		 COALESCE(SUM(wl.quantity), 0)::DECIMAL AS total
		 FROM work_types wt
		 LEFT JOIN work_logs wl ON wl.work_type_id = wt.id
		     AND wl.work_date >= DATE_TRUNC('month', CURRENT_DATE) AND wl.status = 'approved'
		 WHERE wt.active = TRUE
		 GROUP BY wt.id, wt.name, wt.unit, wt.icon, wt.slug, wt.sort_order
		 ORDER BY wt.sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DashboardOutput{}
	for rows.Next() {
		var d DashboardOutput
		if err := rows.Scan(&d.Name, &d.Unit, &d.Icon, &d.Slug, &d.Total); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}

type DashboardCounts struct {
	PendingReviews   int64 `json:"pending_reviews"`
	ContentReview    int64 `json:"content_review"`
	OverdueTasks     int64 `json:"overdue_tasks"`
	RawInbox         int64 `json:"raw_inbox"`
	OpenTasks        int64 `json:"open_tasks"`
	DoneTasks        int64 `json:"done_tasks"`
	RunningCampaigns int64 `json:"running_campaigns"`
}

func (q *Queries) GetDashboardCounts(ctx context.Context) (DashboardCounts, error) {
	var d DashboardCounts
	err := q.pool.QueryRow(ctx,
		`SELECT
		 (SELECT COUNT(*) FROM work_logs WHERE status = 'submitted'),
		 (SELECT COUNT(*) FROM content WHERE status = 'review' AND deleted_at IS NULL),
		 (SELECT COUNT(*) FROM tasks WHERE due_date < CURRENT_DATE AND status NOT IN ('done', 'cancelled') AND deleted_at IS NULL),
		 (SELECT COUNT(*) FROM inbox_items WHERE status = 'raw'),
		 (SELECT COUNT(*) FROM tasks WHERE status NOT IN ('done', 'cancelled') AND deleted_at IS NULL),
		 (SELECT COUNT(*) FROM tasks WHERE status = 'done' AND deleted_at IS NULL),
		 (SELECT COUNT(*) FROM campaigns WHERE status = 'running' AND deleted_at IS NULL)`).
		Scan(&d.PendingReviews, &d.ContentReview, &d.OverdueTasks, &d.RawInbox, &d.OpenTasks, &d.DoneTasks, &d.RunningCampaigns)
	return d, err
}
