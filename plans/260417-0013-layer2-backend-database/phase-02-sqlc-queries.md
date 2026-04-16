# Phase 2: sqlc Queries + Generated Code

## Status: COMPLETE

## Overview
- **Priority:** P0
- **Effort:** 5h (Completed)
- Write SQL queries for all tables, generate type-safe Go code with sqlc

## Completion Summary
All SQL queries written for 12 tables. Generated Go code (models + query functions) hand-written to equivalent standards. Database connection pool configured with pgxpool. All code compiles successfully.

## Implementation Steps

### 1. Create sqlc config
`internal/db/sqlc.yaml`:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query"
    schema: "migrations"
    gen:
      go:
        package: "db"
        out: "generated"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
```

### 2. Write query files (per table)

**users.sql:**
```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY full_name LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, role, status)
VALUES ($1, $2, $3, $4, 'active') RETURNING *;

-- name: UpdateUser :one
UPDATE users SET full_name=$2, role=$3, phone=$4, status=$5
WHERE id=$1 RETURNING *;
```

Similar patterns for:
- **companies.sql** — CRUD + health update + list by status
- **user_company_assignments.sql** — assign user, list by user/company
- **inbox_items.sql** — CRUD + list by status + triage + auto-archive
- **tasks.sql** — CRUD + list by status/assignee/company + kanban counts
- **content.sql** — CRUD + list by status/company + pipeline stats
- **campaigns.sql** — CRUD + progress query + list by company
- **work_types.sql** — list active
- **work_logs.sql** — CRUD + submit/approve/reject + monthly stats
- **knowledge_items.sql** — CRUD + search + list by category
- **objectives.sql** — CRUD + list by company/quarter
- **dashboard.sql** — aggregate queries for dashboard stats

### 3. Generate Go code
```bash
sqlc generate
```

### 4. Create DB connection pool
`internal/db/pool.go`:
```go
func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
    config, _ := pgxpool.ParseConfig(dbURL)
    config.MaxConns = 10
    config.MinConns = 2
    return pgxpool.NewWithConfig(ctx, config)
}
```

### 5. Wire into main.go
```go
pool, err := db.NewPool(ctx, cfg.DatabaseURL)
queries := generated.New(pool)
// Pass queries to handlers
```

## Files Created
- `internal/db/sqlc.yaml` — sqlc config pointing to migrations + queries
- `internal/db/query/*.sql` (12 query files)
  - users.sql (GetByID, GetByEmail, List, Create, Update)
  - companies.sql (CRUD + list by status + health update)
  - user_company_assignments.sql (assign, list by user/company)
  - inbox_items.sql (CRUD + list by status + auto-archive)
  - tasks.sql (CRUD + list by status/assignee + kanban counts)
  - content.sql (CRUD + list by status + pipeline stats)
  - campaigns.sql (CRUD + progress queries)
  - work_types.sql (list active)
  - work_logs.sql (CRUD + submit/approve/reject + monthly stats)
  - knowledge_items.sql (CRUD + search FTS)
  - objectives.sql (CRUD + list by company/quarter)
  - dashboard.sql (aggregate stats queries)
  - bookmarks.sql (CRUD + list by tags)
- `internal/db/pool.go` — pgxpool connection factory with max/min conn config
- `internal/db/generated/*.go` — Hand-written models + query functions (equivalent to sqlc output)

## Files Modified
- (No existing files modified)

## Completed Criteria
✓ All query functions compile without errors
✓ Connection pool successfully connects to remote DB
✓ GetUserByEmail(ctx, "owner@binhvuong.vn") returns seeded user
✓ Type-safe query results with JSON struct tags
✓ All CRUD operations wired for handlers
