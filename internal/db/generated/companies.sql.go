package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *Queries) GetCompanyByID(ctx context.Context, id pgtype.UUID) (Company, error) {
	var c Company
	err := q.pool.QueryRow(ctx,
		`SELECT id, name, short_code, slug, logo_url, industry, my_role, scope, status, health,
		 primary_contact_name, primary_contact_phone, primary_contact_zalo, primary_contact_email,
		 start_date, end_date, description, internal_notes, notion_page_id, synced_at, sync_status, sync_error,
		 created_at, updated_at, created_by, deleted_at
		 FROM companies WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&c.ID, &c.Name, &c.ShortCode, &c.Slug, &c.LogoURL, &c.Industry, &c.MyRole, &c.Scope, &c.Status, &c.Health,
		&c.PrimaryContactName, &c.PrimaryContactPhone, &c.PrimaryContactZalo, &c.PrimaryContactEmail,
		&c.StartDate, &c.EndDate, &c.Description, &c.InternalNotes, &c.NotionPageID, &c.SyncedAt, &c.SyncStatus, &c.SyncError,
		&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.DeletedAt)
	return c, err
}

func scanCompany(scan func(dest ...interface{}) error) (Company, error) {
	var c Company
	err := scan(
		&c.ID, &c.Name, &c.ShortCode, &c.Slug, &c.LogoURL, &c.Industry, &c.MyRole, &c.Scope, &c.Status, &c.Health,
		&c.PrimaryContactName, &c.PrimaryContactPhone, &c.PrimaryContactZalo, &c.PrimaryContactEmail,
		&c.StartDate, &c.EndDate, &c.Description, &c.InternalNotes, &c.NotionPageID, &c.SyncedAt, &c.SyncStatus, &c.SyncError,
		&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.DeletedAt)
	return c, err
}

const companyCols = `id, name, short_code, slug, logo_url, industry, my_role, scope, status, health,
	primary_contact_name, primary_contact_phone, primary_contact_zalo, primary_contact_email,
	start_date, end_date, description, internal_notes, notion_page_id, synced_at, sync_status, sync_error,
	created_at, updated_at, created_by, deleted_at`

func (q *Queries) ListCompanies(ctx context.Context, limit, offset int32) ([]Company, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT `+companyCols+` FROM companies WHERE deleted_at IS NULL ORDER BY name LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Company{}
	for rows.Next() {
		c, err := scanCompany(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (q *Queries) ListCompaniesByStatus(ctx context.Context, status string) ([]Company, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT `+companyCols+` FROM companies WHERE status = $1 AND deleted_at IS NULL ORDER BY name`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Company{}
	for rows.Next() {
		c, err := scanCompany(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (q *Queries) CountCompanies(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM companies WHERE deleted_at IS NULL").Scan(&count)
	return count, err
}

type CreateCompanyParams struct {
	Name        string         `json:"name"`
	ShortCode   sql.NullString `json:"short_code"`
	Slug        sql.NullString `json:"slug"`
	Industry    sql.NullString `json:"industry"`
	MyRole      string         `json:"my_role"`
	Scope       []string       `json:"scope"`
	Status      string         `json:"status"`
	Health      sql.NullString `json:"health"`
	Description sql.NullString `json:"description"`
	CreatedBy   pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error) {
	row := q.pool.QueryRow(ctx,
		`INSERT INTO companies (name, short_code, slug, industry, my_role, scope, status, health, description, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING `+companyCols,
		arg.Name, arg.ShortCode, arg.Slug, arg.Industry, arg.MyRole, arg.Scope, arg.Status, arg.Health, arg.Description, arg.CreatedBy)
	return scanCompany(row.Scan)
}

type UpdateCompanyParams struct {
	ID                  pgtype.UUID    `json:"id"`
	Name                string         `json:"name"`
	ShortCode           sql.NullString `json:"short_code"`
	Industry            sql.NullString `json:"industry"`
	MyRole              string         `json:"my_role"`
	Scope               []string       `json:"scope"`
	Status              string         `json:"status"`
	Health              sql.NullString `json:"health"`
	Description         sql.NullString `json:"description"`
	PrimaryContactName  sql.NullString `json:"primary_contact_name"`
	PrimaryContactPhone sql.NullString `json:"primary_contact_phone"`
	PrimaryContactEmail sql.NullString `json:"primary_contact_email"`
}

func (q *Queries) UpdateCompany(ctx context.Context, arg UpdateCompanyParams) (Company, error) {
	row := q.pool.QueryRow(ctx,
		`UPDATE companies SET name=$2, short_code=$3, industry=$4, my_role=$5,
		 scope=$6, status=$7, health=$8, description=$9,
		 primary_contact_name=$10, primary_contact_phone=$11, primary_contact_email=$12
		 WHERE id=$1 AND deleted_at IS NULL RETURNING `+companyCols,
		arg.ID, arg.Name, arg.ShortCode, arg.Industry, arg.MyRole, arg.Scope, arg.Status, arg.Health,
		arg.Description, arg.PrimaryContactName, arg.PrimaryContactPhone, arg.PrimaryContactEmail)
	return scanCompany(row.Scan)
}

func (q *Queries) UpdateCompanyHealth(ctx context.Context, id pgtype.UUID, health string) error {
	return q.exec(ctx, "UPDATE companies SET health = $2 WHERE id = $1", id, health)
}

type CompanyTaskStats struct {
	CompanyID  pgtype.UUID `json:"company_id"`
	OpenTasks  int64       `json:"open_tasks"`
	DoneTasks  int64       `json:"done_tasks"`
	TotalTasks int64       `json:"total_tasks"`
}

func (q *Queries) GetCompanyTaskStats(ctx context.Context) ([]CompanyTaskStats, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT company_id,
		 COUNT(*) FILTER (WHERE status NOT IN ('done','cancelled')) AS open_tasks,
		 COUNT(*) FILTER (WHERE status = 'done') AS done_tasks,
		 COUNT(*) AS total_tasks
		 FROM tasks WHERE deleted_at IS NULL AND company_id IS NOT NULL
		 GROUP BY company_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CompanyTaskStats{}
	for rows.Next() {
		var s CompanyTaskStats
		if err := rows.Scan(&s.CompanyID, &s.OpenTasks, &s.DoneTasks, &s.TotalTasks); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

func (q *Queries) SoftDeleteCompany(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE companies SET deleted_at = NOW() WHERE id = $1", id)
}
