package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const taskCols = `id, title, description, category, group_name, company_id, assignee_id, objective_id,
	content_id, campaign_id, status, priority, due_date, due_date_end, started_at, completed_at,
	attachments, notion_page_id, synced_at, sync_status, sync_error, created_at, updated_at, created_by, deleted_at`

func scanTask(scan func(dest ...interface{}) error) (Task, error) {
	var t Task
	err := scan(&t.ID, &t.Title, &t.Description, &t.Category, &t.GroupName, &t.CompanyID, &t.AssigneeID,
		&t.ObjectiveID, &t.ContentID, &t.CampaignID, &t.Status, &t.Priority, &t.DueDate, &t.DueDateEnd,
		&t.StartedAt, &t.CompletedAt, &t.Attachments, &t.NotionPageID, &t.SyncedAt, &t.SyncStatus,
		&t.SyncError, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.DeletedAt)
	return t, err
}

func (q *Queries) scanTasks(ctx context.Context, sql string, args ...interface{}) ([]Task, error) {
	rows, err := q.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Task{}
	for rows.Next() {
		t, err := scanTask(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func (q *Queries) GetTaskByID(ctx context.Context, id pgtype.UUID) (Task, error) {
	return scanTask(q.pool.QueryRow(ctx, `SELECT `+taskCols+` FROM tasks WHERE id=$1 AND deleted_at IS NULL`, id).Scan)
}

func (q *Queries) ListTasks(ctx context.Context, limit, offset int32) ([]Task, error) {
	return q.scanTasks(ctx, `SELECT `+taskCols+` FROM tasks WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
}

func (q *Queries) ListTasksByStatus(ctx context.Context, status string) ([]Task, error) {
	return q.scanTasks(ctx, `SELECT `+taskCols+` FROM tasks WHERE status=$1 AND deleted_at IS NULL ORDER BY priority DESC, due_date ASC NULLS LAST`, status)
}

func (q *Queries) ListTasksByCompany(ctx context.Context, companyID pgtype.UUID, limit, offset int32) ([]Task, error) {
	return q.scanTasks(ctx, `SELECT `+taskCols+` FROM tasks WHERE company_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, companyID, limit, offset)
}

func (q *Queries) ListTasksByAssignee(ctx context.Context, assigneeID pgtype.UUID) ([]Task, error) {
	return q.scanTasks(ctx, `SELECT `+taskCols+` FROM tasks WHERE assignee_id=$1 AND deleted_at IS NULL ORDER BY due_date ASC NULLS LAST`, assigneeID)
}

func (q *Queries) ListTasksDueToday(ctx context.Context) ([]Task, error) {
	return q.scanTasks(ctx, `SELECT `+taskCols+` FROM tasks WHERE due_date <= CURRENT_DATE AND status NOT IN ('done','cancelled') AND deleted_at IS NULL ORDER BY priority DESC, due_date ASC`)
}

type TaskStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func (q *Queries) CountTasksByStatus(ctx context.Context) ([]TaskStatusCount, error) {
	rows, err := q.pool.Query(ctx, "SELECT status, COUNT(*) AS count FROM tasks WHERE deleted_at IS NULL GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TaskStatusCount{}
	for rows.Next() {
		var sc TaskStatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		items = append(items, sc)
	}
	return items, rows.Err()
}

func (q *Queries) CountOverdueTasks(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM tasks WHERE due_date < CURRENT_DATE AND status NOT IN ('done','cancelled') AND deleted_at IS NULL").Scan(&count)
	return count, err
}

type CreateTaskParams struct {
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Category    sql.NullString `json:"category"`
	GroupName   sql.NullString `json:"group_name"`
	CompanyID   pgtype.UUID    `json:"company_id"`
	AssigneeID  pgtype.UUID    `json:"assignee_id"`
	CampaignID  pgtype.UUID    `json:"campaign_id"`
	Status      string         `json:"status"`
	Priority    string         `json:"priority"`
	DueDate     pgtype.Date    `json:"due_date"`
	CreatedBy   pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	return scanTask(q.pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, category, group_name, company_id, assignee_id, campaign_id, status, priority, due_date, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING `+taskCols,
		arg.Title, arg.Description, arg.Category, arg.GroupName, arg.CompanyID, arg.AssigneeID,
		arg.CampaignID, arg.Status, arg.Priority, arg.DueDate, arg.CreatedBy).Scan)
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, id pgtype.UUID, status string) (Task, error) {
	return scanTask(q.pool.QueryRow(ctx,
		`UPDATE tasks SET status=$2,
		 started_at = CASE WHEN $2='in_progress' AND started_at IS NULL THEN NOW() ELSE started_at END,
		 completed_at = CASE WHEN $2='done' THEN NOW() ELSE completed_at END
		 WHERE id=$1 AND deleted_at IS NULL RETURNING `+taskCols, id, status).Scan)
}

func (q *Queries) SoftDeleteTask(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE tasks SET deleted_at=NOW() WHERE id=$1", id)
}
