package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const knowledgeCols = `id, title, description, body, category, topics, quality_rating, scope,
	company_id, visible_to_companies, format, source_url, attachments, status,
	notion_page_id, synced_at, sync_status, sync_error, created_at, updated_at, created_by, deleted_at`

func scanKnowledge(scan func(dest ...interface{}) error) (KnowledgeItem, error) {
	var k KnowledgeItem
	err := scan(&k.ID, &k.Title, &k.Description, &k.Body, &k.Category, &k.Topics, &k.QualityRating,
		&k.Scope, &k.CompanyID, &k.VisibleToCompanies, &k.Format, &k.SourceURL, &k.Attachments, &k.Status,
		&k.NotionPageID, &k.SyncedAt, &k.SyncStatus, &k.SyncError, &k.CreatedAt, &k.UpdatedAt, &k.CreatedBy, &k.DeletedAt)
	return k, err
}

func (q *Queries) scanKnowledgeItems(ctx context.Context, sql string, args ...interface{}) ([]KnowledgeItem, error) {
	rows, err := q.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []KnowledgeItem{}
	for rows.Next() {
		k, err := scanKnowledge(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, k)
	}
	return items, rows.Err()
}

func (q *Queries) GetKnowledgeItemByID(ctx context.Context, id pgtype.UUID) (KnowledgeItem, error) {
	return scanKnowledge(q.pool.QueryRow(ctx, `SELECT `+knowledgeCols+` FROM knowledge_items WHERE id=$1 AND deleted_at IS NULL`, id).Scan)
}

func (q *Queries) ListKnowledgeItems(ctx context.Context, limit, offset int32) ([]KnowledgeItem, error) {
	return q.scanKnowledgeItems(ctx, `SELECT `+knowledgeCols+` FROM knowledge_items WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
}

func (q *Queries) ListKnowledgeItemsByCategory(ctx context.Context, category string, limit, offset int32) ([]KnowledgeItem, error) {
	return q.scanKnowledgeItems(ctx, `SELECT `+knowledgeCols+` FROM knowledge_items WHERE category=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, category, limit, offset)
}

func (q *Queries) CountKnowledgeItems(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM knowledge_items WHERE deleted_at IS NULL").Scan(&count)
	return count, err
}

func (q *Queries) SearchKnowledgeItems(ctx context.Context, query string, limit, offset int32) ([]KnowledgeItem, error) {
	return q.scanKnowledgeItems(ctx,
		`SELECT `+knowledgeCols+` FROM knowledge_items WHERE deleted_at IS NULL
		 AND to_tsvector('simple', unaccent(title || ' ' || COALESCE(description, ''))) @@ plainto_tsquery('simple', unaccent($1))
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, query, limit, offset)
}

type CreateKnowledgeItemParams struct {
	Title         string         `json:"title"`
	Description   sql.NullString `json:"description"`
	Body          sql.NullString `json:"body"`
	Category      string         `json:"category"`
	Topics        []string       `json:"topics"`
	QualityRating pgtype.Int4    `json:"quality_rating"`
	Scope         string         `json:"scope"`
	CompanyID     pgtype.UUID    `json:"company_id"`
	Format        sql.NullString `json:"format"`
	SourceURL     sql.NullString `json:"source_url"`
	CreatedBy     pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateKnowledgeItem(ctx context.Context, arg CreateKnowledgeItemParams) (KnowledgeItem, error) {
	return scanKnowledge(q.pool.QueryRow(ctx,
		`INSERT INTO knowledge_items (title, description, body, category, topics, quality_rating, scope, company_id, format, source_url, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING `+knowledgeCols,
		arg.Title, arg.Description, arg.Body, arg.Category, arg.Topics, arg.QualityRating,
		arg.Scope, arg.CompanyID, arg.Format, arg.SourceURL, arg.CreatedBy).Scan)
}

func (q *Queries) SoftDeleteKnowledgeItem(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE knowledge_items SET deleted_at=NOW() WHERE id=$1", id)
}
