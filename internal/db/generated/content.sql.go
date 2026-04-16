package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const contentCols = `id, title, content_type, platforms, topics, company_id, author_id, campaign_id, status,
	publish_date, published_url, source_file_url, attachments, reach, engagement, visible_to_companies,
	notes, review_notes, notion_page_id, synced_at, sync_status, sync_error,
	created_at, updated_at, created_by, deleted_at, engagement_rate`

func scanContent(scan func(dest ...interface{}) error) (Content, error) {
	var c Content
	err := scan(&c.ID, &c.Title, &c.ContentType, &c.Platforms, &c.Topics, &c.CompanyID, &c.AuthorID,
		&c.CampaignID, &c.Status, &c.PublishDate, &c.PublishedURL, &c.SourceFileURL, &c.Attachments,
		&c.Reach, &c.Engagement, &c.VisibleToCompanies, &c.Notes, &c.ReviewNotes,
		&c.NotionPageID, &c.SyncedAt, &c.SyncStatus, &c.SyncError,
		&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.DeletedAt, &c.EngagementRate)
	return c, err
}

func (q *Queries) scanContents(ctx context.Context, sql string, args ...interface{}) ([]Content, error) {
	rows, err := q.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Content{}
	for rows.Next() {
		c, err := scanContent(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (q *Queries) GetContentByID(ctx context.Context, id pgtype.UUID) (Content, error) {
	return scanContent(q.pool.QueryRow(ctx, `SELECT `+contentCols+` FROM content WHERE id=$1 AND deleted_at IS NULL`, id).Scan)
}

func (q *Queries) ListContent(ctx context.Context, limit, offset int32) ([]Content, error) {
	return q.scanContents(ctx, `SELECT `+contentCols+` FROM content WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
}

func (q *Queries) ListContentByStatus(ctx context.Context, status string, limit, offset int32) ([]Content, error) {
	return q.scanContents(ctx, `SELECT `+contentCols+` FROM content WHERE status=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
}

func (q *Queries) ListContentByCompany(ctx context.Context, companyID pgtype.UUID, limit, offset int32) ([]Content, error) {
	return q.scanContents(ctx, `SELECT `+contentCols+` FROM content WHERE company_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, companyID, limit, offset)
}

func (q *Queries) CountContent(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM content WHERE deleted_at IS NULL").Scan(&count)
	return count, err
}

type ContentStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func (q *Queries) CountContentByStatus(ctx context.Context) ([]ContentStatusCount, error) {
	rows, err := q.pool.Query(ctx, "SELECT status, COUNT(*) AS count FROM content WHERE deleted_at IS NULL GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ContentStatusCount{}
	for rows.Next() {
		var sc ContentStatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		items = append(items, sc)
	}
	return items, rows.Err()
}

type CreateContentParams struct {
	Title       string         `json:"title"`
	ContentType string         `json:"content_type"`
	Platforms   []string       `json:"platforms"`
	Topics      []string       `json:"topics"`
	CompanyID   pgtype.UUID    `json:"company_id"`
	AuthorID    pgtype.UUID    `json:"author_id"`
	CampaignID  pgtype.UUID    `json:"campaign_id"`
	Status      string         `json:"status"`
	CreatedBy   pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateContent(ctx context.Context, arg CreateContentParams) (Content, error) {
	return scanContent(q.pool.QueryRow(ctx,
		`INSERT INTO content (title, content_type, platforms, topics, company_id, author_id, campaign_id, status, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING `+contentCols,
		arg.Title, arg.ContentType, arg.Platforms, arg.Topics, arg.CompanyID, arg.AuthorID, arg.CampaignID, arg.Status, arg.CreatedBy).Scan)
}

func (q *Queries) ReviewContent(ctx context.Context, id pgtype.UUID, status string, reviewNotes sql.NullString) (Content, error) {
	return scanContent(q.pool.QueryRow(ctx,
		`UPDATE content SET status=$2, review_notes=$3 WHERE id=$1 AND deleted_at IS NULL RETURNING `+contentCols,
		id, status, reviewNotes).Scan)
}

func (q *Queries) PublishContent(ctx context.Context, id pgtype.UUID, publishDate pgtype.Date, publishedURL sql.NullString) (Content, error) {
	return scanContent(q.pool.QueryRow(ctx,
		`UPDATE content SET status='published', publish_date=$2, published_url=$3 WHERE id=$1 AND deleted_at IS NULL RETURNING `+contentCols,
		id, publishDate, publishedURL).Scan)
}

func (q *Queries) SoftDeleteContent(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE content SET deleted_at=NOW() WHERE id=$1", id)
}
