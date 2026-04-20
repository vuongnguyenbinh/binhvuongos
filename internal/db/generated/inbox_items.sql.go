package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const inboxCols = `id, content, url, source, item_type, status, destination, company_id, submitted_by,
	attachments, telegram_message_id, telegram_chat_id, triage_notes, processed_at,
	converted_to_type, converted_to_id, external_ref, created_at, updated_at`

func scanInbox(scan func(dest ...interface{}) error) (InboxItem, error) {
	var i InboxItem
	err := scan(&i.ID, &i.Content, &i.URL, &i.Source, &i.ItemType, &i.Status, &i.Destination,
		&i.CompanyID, &i.SubmittedBy, &i.Attachments, &i.TelegramMessageID, &i.TelegramChatID,
		&i.TriageNotes, &i.ProcessedAt, &i.ConvertedToType, &i.ConvertedToID, &i.ExternalRef,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

func (q *Queries) GetInboxItemByID(ctx context.Context, id pgtype.UUID) (InboxItem, error) {
	return scanInbox(q.pool.QueryRow(ctx, `SELECT `+inboxCols+` FROM inbox_items WHERE id = $1`, id).Scan)
}

func (q *Queries) ListInboxItems(ctx context.Context, limit, offset int32) ([]InboxItem, error) {
	rows, err := q.pool.Query(ctx, `SELECT `+inboxCols+` FROM inbox_items ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []InboxItem{}
	for rows.Next() {
		i, err := scanInbox(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

func (q *Queries) ListInboxItemsByStatus(ctx context.Context, status string, limit, offset int32) ([]InboxItem, error) {
	rows, err := q.pool.Query(ctx, `SELECT `+inboxCols+` FROM inbox_items WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []InboxItem{}
	for rows.Next() {
		i, err := scanInbox(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

func (q *Queries) CountInboxItems(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM inbox_items").Scan(&count)
	return count, err
}

func (q *Queries) CountInboxItemsByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM inbox_items WHERE status = $1", status).Scan(&count)
	return count, err
}

type CreateInboxItemParams struct {
	Content     string         `json:"content"`
	URL         sql.NullString `json:"url"`
	Source      sql.NullString `json:"source"`
	ItemType    sql.NullString `json:"item_type"`
	CompanyID   pgtype.UUID    `json:"company_id"`
	SubmittedBy pgtype.UUID    `json:"submitted_by"`
	Attachments []byte         `json:"attachments"`
	ExternalRef sql.NullString `json:"external_ref"`
}

func (q *Queries) CreateInboxItem(ctx context.Context, arg CreateInboxItemParams) (InboxItem, error) {
	// attachments defaults to '[]' when nil to satisfy JSONB column
	att := arg.Attachments
	if att == nil {
		att = []byte("[]")
	}
	return scanInbox(q.pool.QueryRow(ctx,
		`INSERT INTO inbox_items (content, url, source, item_type, company_id, submitted_by, attachments, external_ref)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING `+inboxCols,
		arg.Content, arg.URL, arg.Source, arg.ItemType, arg.CompanyID, arg.SubmittedBy, att, arg.ExternalRef).Scan)
}

// GetInboxByExternalRef returns an inbox item matching (source, external_ref) if any.
// Used for idempotency — retried webhook deliveries won't create duplicates.
type GetInboxByExternalRefParams struct {
	Source      sql.NullString `json:"source"`
	ExternalRef sql.NullString `json:"external_ref"`
}

func (q *Queries) GetInboxByExternalRef(ctx context.Context, arg GetInboxByExternalRefParams) (InboxItem, error) {
	return scanInbox(q.pool.QueryRow(ctx,
		`SELECT `+inboxCols+` FROM inbox_items WHERE source = $1 AND external_ref = $2 LIMIT 1`,
		arg.Source, arg.ExternalRef).Scan)
}

// UpsertInboxItemByExternalRef inserts a new inbox row or returns the existing one when
// (source, external_ref) collides. Atomic — no race between concurrent retries.
// Second return value `inserted` is true for a new row, false when an existing row was returned.
func (q *Queries) UpsertInboxItemByExternalRef(ctx context.Context, arg CreateInboxItemParams) (InboxItem, bool, error) {
	att := arg.Attachments
	if att == nil {
		att = []byte("[]")
	}
	// Try insert; on unique conflict do nothing (returns no rows)
	item, err := scanInbox(q.pool.QueryRow(ctx,
		`INSERT INTO inbox_items (content, url, source, item_type, company_id, submitted_by, attachments, external_ref)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (source, external_ref) WHERE external_ref IS NOT NULL DO NOTHING
		 RETURNING `+inboxCols,
		arg.Content, arg.URL, arg.Source, arg.ItemType, arg.CompanyID, arg.SubmittedBy, att, arg.ExternalRef).Scan)
	if err == nil {
		return item, true, nil
	}
	// Conflict → fetch existing row (pgx returns ErrNoRows when no row returned)
	existing, gerr := q.GetInboxByExternalRef(ctx, GetInboxByExternalRefParams{
		Source: arg.Source, ExternalRef: arg.ExternalRef,
	})
	if gerr == nil {
		return existing, false, nil
	}
	// Real insert error (not a conflict-then-miss) — propagate the original
	return InboxItem{}, false, err
}

type UpdateInboxItemParams struct {
	ID       pgtype.UUID    `json:"id"`
	Content  string         `json:"content"`
	URL      sql.NullString `json:"url"`
	ItemType sql.NullString `json:"item_type"`
}

func (q *Queries) UpdateInboxItem(ctx context.Context, arg UpdateInboxItemParams) (InboxItem, error) {
	return scanInbox(q.pool.QueryRow(ctx,
		`UPDATE inbox_items SET content=$2, url=$3, item_type=$4 WHERE id=$1 RETURNING `+inboxCols,
		arg.ID, arg.Content, arg.URL, arg.ItemType).Scan)
}

type TriageInboxItemParams struct {
	ID              pgtype.UUID    `json:"id"`
	Destination     sql.NullString `json:"destination"`
	TriageNotes     sql.NullString `json:"triage_notes"`
	ConvertedToType sql.NullString `json:"converted_to_type"`
	ConvertedToID   pgtype.UUID    `json:"converted_to_id"`
}

func (q *Queries) TriageInboxItem(ctx context.Context, arg TriageInboxItemParams) (InboxItem, error) {
	return scanInbox(q.pool.QueryRow(ctx,
		`UPDATE inbox_items SET status='done', destination=$2, triage_notes=$3,
		 converted_to_type=$4, converted_to_id=$5, processed_at=NOW()
		 WHERE id=$1 RETURNING `+inboxCols,
		arg.ID, arg.Destination, arg.TriageNotes, arg.ConvertedToType, arg.ConvertedToID).Scan)
}

func (q *Queries) ArchiveInboxItem(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE inbox_items SET status='archived' WHERE id=$1", id)
}

func (q *Queries) ArchiveOldInboxItems(ctx context.Context) error {
	return q.exec(ctx, "UPDATE inbox_items SET status='archived' WHERE status='raw' AND created_at < NOW() - INTERVAL '7 days'")
}
