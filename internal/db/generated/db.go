// Code generated manually (sqlc not available locally, runs in Docker build).
package generated

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Queries {
	return &Queries{pool: pool}
}

// queryRow executes a query that returns a single row
func (q *Queries) queryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return q.pool.QueryRow(ctx, sql, args...)
}

// query executes a query that returns multiple rows
func (q *Queries) query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return q.pool.Query(ctx, sql, args...)
}

// exec executes a query that doesn't return rows
func (q *Queries) exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := q.pool.Exec(ctx, sql, args...)
	return err
}
