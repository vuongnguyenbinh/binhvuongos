package generated

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Notification struct {
	ID        pgtype.UUID        `json:"id"`
	UserID    pgtype.UUID        `json:"user_id"`
	Title     string             `json:"title"`
	Body      *string            `json:"body"`
	Link      *string            `json:"link"`
	ReadAt    pgtype.Timestamptz `json:"read_at"`
	CreatedAt time.Time          `json:"created_at"`
}

func (q *Queries) ListNotifications(ctx context.Context, userID pgtype.UUID, limit int32) ([]Notification, error) {
	rows, err := q.pool.Query(ctx,
		"SELECT id, user_id, title, body, link, read_at, created_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2",
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.Link, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return items, rows.Err()
}

func (q *Queries) CountUnreadNotifications(ctx context.Context, userID pgtype.UUID) (int64, error) {
	var count int64
	err := q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND read_at IS NULL", userID).Scan(&count)
	return count, err
}

func (q *Queries) MarkNotificationRead(ctx context.Context, id pgtype.UUID) error {
	_, err := q.pool.Exec(ctx, "UPDATE notifications SET read_at=NOW() WHERE id=$1", id)
	return err
}

func (q *Queries) MarkAllNotificationsRead(ctx context.Context, userID pgtype.UUID) error {
	_, err := q.pool.Exec(ctx, "UPDATE notifications SET read_at=NOW() WHERE user_id=$1 AND read_at IS NULL", userID)
	return err
}

func (q *Queries) CreateNotification(ctx context.Context, userID pgtype.UUID, title string, body *string, link *string) error {
	_, err := q.pool.Exec(ctx,
		"INSERT INTO notifications (user_id, title, body, link) VALUES ($1, $2, $3, $4)",
		userID, title, body, link)
	return err
}
