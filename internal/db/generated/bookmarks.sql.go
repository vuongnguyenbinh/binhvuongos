package generated

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const bookmarkCols = `id, title, url, description, tags, notes, created_by, created_at, updated_at, deleted_at`

func scanBookmark(scan func(dest ...interface{}) error) (Bookmark, error) {
	var b Bookmark
	err := scan(&b.ID, &b.Title, &b.URL, &b.Description, &b.Tags, &b.Notes, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)
	return b, err
}

func (q *Queries) GetBookmarkByID(ctx context.Context, id pgtype.UUID) (Bookmark, error) {
	return scanBookmark(q.pool.QueryRow(ctx, `SELECT `+bookmarkCols+` FROM bookmarks WHERE id=$1 AND deleted_at IS NULL`, id).Scan)
}

func (q *Queries) ListBookmarks(ctx context.Context, limit, offset int32) ([]Bookmark, error) {
	rows, err := q.pool.Query(ctx, `SELECT `+bookmarkCols+` FROM bookmarks WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Bookmark{}
	for rows.Next() {
		b, err := scanBookmark(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, rows.Err()
}

func (q *Queries) CountBookmarks(ctx context.Context) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookmarks WHERE deleted_at IS NULL").Scan(&count)
	return count, err
}

type CreateBookmarkParams struct {
	Title       string         `json:"title"`
	URL         string         `json:"url"`
	Description sql.NullString `json:"description"`
	Tags        []string       `json:"tags"`
	Notes       sql.NullString `json:"notes"`
	CreatedBy   pgtype.UUID    `json:"created_by"`
}

func (q *Queries) CreateBookmark(ctx context.Context, arg CreateBookmarkParams) (Bookmark, error) {
	return scanBookmark(q.pool.QueryRow(ctx,
		`INSERT INTO bookmarks (title, url, description, tags, notes, created_by) VALUES ($1,$2,$3,$4,$5,$6) RETURNING `+bookmarkCols,
		arg.Title, arg.URL, arg.Description, arg.Tags, arg.Notes, arg.CreatedBy).Scan)
}

func (q *Queries) SoftDeleteBookmark(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "UPDATE bookmarks SET deleted_at=NOW() WHERE id=$1", id)
}
