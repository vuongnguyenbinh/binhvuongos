package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const objectiveCols = `id, title, description, company_id, owner_id, quarter, year, status, progress,
	key_results, notes, notion_page_id, synced_at, sync_status, sync_error,
	created_at, updated_at, created_by, deleted_at`

func scanObjective(scan func(dest ...interface{}) error) (Objective, error) {
	var o Objective
	err := scan(&o.ID, &o.Title, &o.Description, &o.CompanyID, &o.OwnerID, &o.Quarter, &o.Year,
		&o.Status, &o.Progress, &o.KeyResults, &o.Notes, &o.NotionPageID, &o.SyncedAt, &o.SyncStatus,
		&o.SyncError, &o.CreatedAt, &o.UpdatedAt, &o.CreatedBy, &o.DeletedAt)
	return o, err
}

func (q *Queries) GetObjectiveByID(ctx context.Context, id pgtype.UUID) (Objective, error) {
	return scanObjective(q.pool.QueryRow(ctx, `SELECT `+objectiveCols+` FROM objectives WHERE id=$1 AND deleted_at IS NULL`, id).Scan)
}

func (q *Queries) ListObjectivesByCompany(ctx context.Context, companyID pgtype.UUID) ([]Objective, error) {
	rows, err := q.pool.Query(ctx, `SELECT `+objectiveCols+` FROM objectives WHERE company_id=$1 AND deleted_at IS NULL ORDER BY quarter DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Objective{}
	for rows.Next() {
		o, err := scanObjective(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	return items, rows.Err()
}

type CreateObjectiveParams struct {
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	CompanyID   pgtype.UUID    `json:"company_id"`
	OwnerID     pgtype.UUID    `json:"owner_id"`
	Quarter     string         `json:"quarter"`
	Year        pgtype.Int4    `json:"year"`
	KeyResults  []byte         `json:"key_results"`
	CreatedBy   pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateObjective(ctx context.Context, arg CreateObjectiveParams) (Objective, error) {
	return scanObjective(q.pool.QueryRow(ctx,
		`INSERT INTO objectives (title, description, company_id, owner_id, quarter, year, key_results, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+objectiveCols,
		arg.Title, arg.Description, arg.CompanyID, arg.OwnerID, arg.Quarter, arg.Year, arg.KeyResults, arg.CreatedBy).Scan)
}
