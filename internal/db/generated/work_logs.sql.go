package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const workLogCols = `id, work_date, user_id, company_id, work_type_id, campaign_id, quantity,
	sheet_url, evidence_url, screenshots, notes, admin_notes, status,
	submitted_at, reviewed_at, reviewed_by, notion_page_id, synced_at, sync_status, sync_error,
	created_at, updated_at`

func scanWorkLog(scan func(dest ...interface{}) error) (WorkLog, error) {
	var w WorkLog
	err := scan(&w.ID, &w.WorkDate, &w.UserID, &w.CompanyID, &w.WorkTypeID, &w.CampaignID, &w.Quantity,
		&w.SheetURL, &w.EvidenceURL, &w.Screenshots, &w.Notes, &w.AdminNotes, &w.Status,
		&w.SubmittedAt, &w.ReviewedAt, &w.ReviewedBy, &w.NotionPageID, &w.SyncedAt, &w.SyncStatus, &w.SyncError,
		&w.CreatedAt, &w.UpdatedAt)
	return w, err
}

func (q *Queries) scanWorkLogs(ctx context.Context, sql string, args ...interface{}) ([]WorkLog, error) {
	rows, err := q.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []WorkLog{}
	for rows.Next() {
		w, err := scanWorkLog(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, w)
	}
	return items, rows.Err()
}

func (q *Queries) GetWorkLogByID(ctx context.Context, id pgtype.UUID) (WorkLog, error) {
	return scanWorkLog(q.pool.QueryRow(ctx, `SELECT `+workLogCols+` FROM work_logs WHERE id=$1`, id).Scan)
}

func (q *Queries) ListWorkLogs(ctx context.Context, limit, offset int32) ([]WorkLog, error) {
	return q.scanWorkLogs(ctx, `SELECT `+workLogCols+` FROM work_logs ORDER BY work_date DESC, created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
}

func (q *Queries) ListWorkLogsByStatus(ctx context.Context, status string, limit, offset int32) ([]WorkLog, error) {
	return q.scanWorkLogs(ctx, `SELECT `+workLogCols+` FROM work_logs WHERE status=$1 ORDER BY work_date DESC LIMIT $2 OFFSET $3`, status, limit, offset)
}

func (q *Queries) ListWorkLogsByUser(ctx context.Context, userID pgtype.UUID, limit, offset int32) ([]WorkLog, error) {
	return q.scanWorkLogs(ctx, `SELECT `+workLogCols+` FROM work_logs WHERE user_id=$1 ORDER BY work_date DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
}

func (q *Queries) ListWorkLogsByCompany(ctx context.Context, companyID pgtype.UUID, limit, offset int32) ([]WorkLog, error) {
	return q.scanWorkLogs(ctx, `SELECT `+workLogCols+` FROM work_logs WHERE company_id=$1 ORDER BY work_date DESC LIMIT $2 OFFSET $3`, companyID, limit, offset)
}

func (q *Queries) CountWorkLogs(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM work_logs").Scan(&count)
	return count, err
}

func (q *Queries) CountPendingWorkLogs(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM work_logs WHERE status='submitted'").Scan(&count)
	return count, err
}

type CreateWorkLogParams struct {
	WorkDate    pgtype.Date    `json:"work_date"`
	UserID      pgtype.UUID    `json:"user_id"`
	CompanyID   pgtype.UUID    `json:"company_id"`
	WorkTypeID  pgtype.UUID    `json:"work_type_id"`
	CampaignID  pgtype.UUID    `json:"campaign_id"`
	Quantity    pgtype.Numeric `json:"quantity"`
	SheetURL    sql.NullString `json:"sheet_url"`
	EvidenceURL sql.NullString `json:"evidence_url"`
	Notes       sql.NullString `json:"notes"`
}

func (q *Queries) CreateWorkLog(ctx context.Context, arg CreateWorkLogParams) (WorkLog, error) {
	return scanWorkLog(q.pool.QueryRow(ctx,
		`INSERT INTO work_logs (work_date, user_id, company_id, work_type_id, campaign_id, quantity, sheet_url, evidence_url, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING `+workLogCols,
		arg.WorkDate, arg.UserID, arg.CompanyID, arg.WorkTypeID, arg.CampaignID, arg.Quantity,
		arg.SheetURL, arg.EvidenceURL, arg.Notes).Scan)
}

func (q *Queries) ApproveWorkLog(ctx context.Context, id pgtype.UUID, reviewedBy pgtype.UUID, adminNotes sql.NullString) (WorkLog, error) {
	return scanWorkLog(q.pool.QueryRow(ctx,
		`UPDATE work_logs SET status='approved', reviewed_at=NOW(), reviewed_by=$2, admin_notes=$3 WHERE id=$1 RETURNING `+workLogCols,
		id, reviewedBy, adminNotes).Scan)
}

func (q *Queries) RejectWorkLog(ctx context.Context, id pgtype.UUID, reviewedBy pgtype.UUID, adminNotes sql.NullString) (WorkLog, error) {
	return scanWorkLog(q.pool.QueryRow(ctx,
		`UPDATE work_logs SET status='rejected', reviewed_at=NOW(), reviewed_by=$2, admin_notes=$3 WHERE id=$1 RETURNING `+workLogCols,
		id, reviewedBy, adminNotes).Scan)
}
