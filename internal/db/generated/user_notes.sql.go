package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *Queries) GetUserNote(ctx context.Context, userID pgtype.UUID) (string, error) {
	var content string
	err := q.pool.QueryRow(ctx, "SELECT content FROM user_notes WHERE user_id = $1", userID).Scan(&content)
	return content, err
}

func (q *Queries) UpsertUserNote(ctx context.Context, userID pgtype.UUID, content string) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO user_notes (user_id, content, updated_at) VALUES ($1, $2, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET content = $2, updated_at = NOW()`,
		userID, content)
	return err
}
