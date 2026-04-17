package generated

import (
	"context"
	"time"
)

type Setting struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := q.pool.QueryRow(ctx, "SELECT value FROM settings WHERE key = $1", key).Scan(&value)
	return value, err
}

func (q *Queries) SetSetting(ctx context.Context, key, value string) error {
	_, err := q.pool.Exec(ctx,
		"INSERT INTO settings (key, value, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()",
		key, value)
	return err
}

func (q *Queries) ListSettings(ctx context.Context) ([]Setting, error) {
	rows, err := q.pool.Query(ctx, "SELECT key, value, updated_at FROM settings ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Setting{}
	for rows.Next() {
		var s Setting
		if err := rows.Scan(&s.Key, &s.Value, &s.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}
