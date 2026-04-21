package generated

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateResetToken inserts a single-use password reset token valid for 1 hour.
func (q *Queries) CreateResetToken(ctx context.Context, userID pgtype.UUID, token string) (PasswordResetToken, error) {
	row := q.pool.QueryRow(ctx,
		`INSERT INTO password_reset_tokens (user_id, token, expires_at)
		 VALUES ($1, $2, NOW() + INTERVAL '1 hour')
		 RETURNING id, user_id, token, expires_at, used_at, created_at`,
		userID, token)
	var t PasswordResetToken
	err := row.Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	return t, err
}

// GetValidResetToken returns the token row only when it is unused and unexpired.
func (q *Queries) GetValidResetToken(ctx context.Context, token string) (PasswordResetToken, error) {
	row := q.pool.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at, used_at, created_at
		 FROM password_reset_tokens
		 WHERE token = $1 AND used_at IS NULL AND expires_at > NOW()
		 LIMIT 1`,
		token)
	var t PasswordResetToken
	err := row.Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	return t, err
}

// MarkResetTokenUsed flags a token consumed so it cannot be replayed.
func (q *Queries) MarkResetTokenUsed(ctx context.Context, token string) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE password_reset_tokens SET used_at = $2 WHERE token = $1`,
		token, time.Now())
	return err
}
