package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *Queries) ListActiveWorkTypes(ctx context.Context) ([]WorkType, error) {
	rows, err := q.pool.Query(ctx, "SELECT id, name, slug, unit, icon, color, description, active, sort_order, default_target_per_day, created_at, updated_at FROM work_types WHERE active=TRUE ORDER BY sort_order")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []WorkType{}
	for rows.Next() {
		var w WorkType
		if err := rows.Scan(&w.ID, &w.Name, &w.Slug, &w.Unit, &w.Icon, &w.Color, &w.Description,
			&w.Active, &w.SortOrder, &w.DefaultTargetPerDay, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, w)
	}
	return items, rows.Err()
}

func (q *Queries) GetWorkTypeByID(ctx context.Context, id pgtype.UUID) (WorkType, error) {
	var w WorkType
	err := q.pool.QueryRow(ctx, "SELECT id, name, slug, unit, icon, color, description, active, sort_order, default_target_per_day, created_at, updated_at FROM work_types WHERE id=$1", id).
		Scan(&w.ID, &w.Name, &w.Slug, &w.Unit, &w.Icon, &w.Color, &w.Description,
			&w.Active, &w.SortOrder, &w.DefaultTargetPerDay, &w.CreatedAt, &w.UpdatedAt)
	return w, err
}
