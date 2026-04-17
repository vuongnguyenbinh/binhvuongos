package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const campaignCols = `id, name, description, company_id, owner_id, campaign_type, status,
	start_date, end_date, target_json, budget, budget_spent, notes,
	notion_page_id, synced_at, sync_status, sync_error, created_at, updated_at, created_by, deleted_at`

func scanCampaign(scan func(dest ...interface{}) error) (Campaign, error) {
	var c Campaign
	err := scan(&c.ID, &c.Name, &c.Description, &c.CompanyID, &c.OwnerID, &c.CampaignType, &c.Status,
		&c.StartDate, &c.EndDate, &c.TargetJSON, &c.Budget, &c.BudgetSpent, &c.Notes,
		&c.NotionPageID, &c.SyncedAt, &c.SyncStatus, &c.SyncError, &c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.DeletedAt)
	return c, err
}

func (q *Queries) scanCampaigns(ctx context.Context, sql string, args ...interface{}) ([]Campaign, error) {
	rows, err := q.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Campaign{}
	for rows.Next() {
		c, err := scanCampaign(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (q *Queries) GetCampaignByID(ctx context.Context, id pgtype.UUID) (Campaign, error) {
	return scanCampaign(q.pool.QueryRow(ctx, `SELECT `+campaignCols+` FROM campaigns WHERE id=$1 AND deleted_at IS NULL`, id).Scan)
}

func (q *Queries) ListCampaigns(ctx context.Context, limit, offset int32) ([]Campaign, error) {
	return q.scanCampaigns(ctx, `SELECT `+campaignCols+` FROM campaigns WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
}

func (q *Queries) ListCampaignsByCompany(ctx context.Context, companyID pgtype.UUID) ([]Campaign, error) {
	return q.scanCampaigns(ctx, `SELECT `+campaignCols+` FROM campaigns WHERE company_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC`, companyID)
}

func (q *Queries) ListCampaignsByStatus(ctx context.Context, status string) ([]Campaign, error) {
	return q.scanCampaigns(ctx, `SELECT `+campaignCols+` FROM campaigns WHERE status=$1 AND deleted_at IS NULL ORDER BY start_date DESC`, status)
}

func (q *Queries) CountCampaigns(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM campaigns WHERE deleted_at IS NULL").Scan(&count)
	return count, err
}

type CreateCampaignParams struct {
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description"`
	CompanyID    pgtype.UUID    `json:"company_id"`
	OwnerID      pgtype.UUID    `json:"owner_id"`
	CampaignType sql.NullString `json:"campaign_type"`
	Status       string         `json:"status"`
	StartDate    pgtype.Date    `json:"start_date"`
	EndDate      pgtype.Date    `json:"end_date"`
	TargetJSON   []byte         `json:"target_json"`
	Budget       pgtype.Numeric `json:"budget"`
	CreatedBy    pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateCampaign(ctx context.Context, arg CreateCampaignParams) (Campaign, error) {
	return scanCampaign(q.pool.QueryRow(ctx,
		`INSERT INTO campaigns (name, description, company_id, owner_id, campaign_type, status, start_date, end_date, target_json, budget, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING `+campaignCols,
		arg.Name, arg.Description, arg.CompanyID, arg.OwnerID, arg.CampaignType, arg.Status,
		arg.StartDate, arg.EndDate, arg.TargetJSON, arg.Budget, arg.CreatedBy).Scan)
}

type UpdateCampaignParams struct {
	ID           pgtype.UUID    `json:"id"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description"`
	CampaignType sql.NullString `json:"campaign_type"`
	Status       string         `json:"status"`
	StartDate    pgtype.Date    `json:"start_date"`
	EndDate      pgtype.Date    `json:"end_date"`
}

func (q *Queries) UpdateCampaign(ctx context.Context, arg UpdateCampaignParams) (Campaign, error) {
	return scanCampaign(q.pool.QueryRow(ctx,
		`UPDATE campaigns SET name=$2, description=$3, campaign_type=$4, status=$5, start_date=$6, end_date=$7
		 WHERE id=$1 AND deleted_at IS NULL RETURNING `+campaignCols,
		arg.ID, arg.Name, arg.Description, arg.CampaignType, arg.Status, arg.StartDate, arg.EndDate).Scan)
}

func (q *Queries) SoftDeleteCampaign(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE campaigns SET deleted_at=NOW() WHERE id=$1", id)
}
