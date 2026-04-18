package generated

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CommentWithUser struct {
	ID        pgtype.UUID `json:"id"`
	Body      string      `json:"body"`
	UserID    pgtype.UUID `json:"user_id"`
	FullName  string      `json:"full_name"`
	AvatarURL *string     `json:"avatar_url"`
	CreatedAt time.Time   `json:"created_at"`
}

func (q *Queries) ListCommentsByEntity(ctx context.Context, entityType string, entityID pgtype.UUID) ([]CommentWithUser, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT c.id, c.body, c.user_id, u.full_name, u.avatar_url, c.created_at
		 FROM comments c JOIN users u ON u.id = c.user_id
		 WHERE c.entity_type = $1 AND c.entity_id = $2
		 ORDER BY c.created_at ASC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CommentWithUser{}
	for rows.Next() {
		var c CommentWithUser
		if err := rows.Scan(&c.ID, &c.Body, &c.UserID, &c.FullName, &c.AvatarURL, &c.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (q *Queries) CreateComment(ctx context.Context, entityType string, entityID, userID pgtype.UUID, body string) error {
	_, err := q.pool.Exec(ctx,
		"INSERT INTO comments (entity_type, entity_id, user_id, body) VALUES ($1, $2, $3, $4)",
		entityType, entityID, userID, body)
	return err
}

func (q *Queries) DeleteComment(ctx context.Context, id pgtype.UUID) error {
	_, err := q.pool.Exec(ctx, "DELETE FROM comments WHERE id = $1", id)
	return err
}

func (q *Queries) CountCommentsByEntity(ctx context.Context, entityType string, entityID pgtype.UUID) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM comments WHERE entity_type = $1 AND entity_id = $2",
		entityType, entityID).Scan(&count)
	return count, err
}
