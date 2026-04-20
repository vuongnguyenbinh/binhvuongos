# Phase 1 — Migration + sqlc Query

**Effort:** 20m | **Priority:** P1 | **Status:** pending

## Context
- Brainstorm: [../reports/brainstorm-260420-2131-unified-inbox-webhook.md](../reports/brainstorm-260420-2131-unified-inbox-webhook.md)
- Existing migrations: `internal/db/migrations/000001` → `000021`
- sqlc config: `sqlc.yaml` (check before run generate)

## Overview
Add `external_ref` column + unique index to `inbox_items` for idempotency. Add sqlc query to lookup by `(source, external_ref)`.

## Requirements
- Column `external_ref VARCHAR(200)` nullable
- Unique partial index on `(source, external_ref) WHERE external_ref IS NOT NULL`
- sqlc query returns nullable row (no error if not found)

## Files to create
- `internal/db/migrations/000022_inbox_external_ref.up.sql`
- `internal/db/migrations/000022_inbox_external_ref.down.sql`

## Files to modify
- `internal/db/query/inbox_items.sql` — add `GetInboxByExternalRef`
- `internal/db/generated/inbox_items.sql.go` — auto-regen via `sqlc generate`

## Implementation steps

### 1. Create up migration
```sql
-- 000022_inbox_external_ref.up.sql
ALTER TABLE inbox_items ADD COLUMN external_ref VARCHAR(200);

CREATE UNIQUE INDEX idx_inbox_source_external_ref
  ON inbox_items(source, external_ref)
  WHERE external_ref IS NOT NULL;
```

### 2. Create down migration
```sql
-- 000022_inbox_external_ref.down.sql
DROP INDEX IF EXISTS idx_inbox_source_external_ref;
ALTER TABLE inbox_items DROP COLUMN IF EXISTS external_ref;
```

### 3. Add sqlc query
Append to `internal/db/query/inbox_items.sql`:
```sql
-- name: GetInboxByExternalRef :one
SELECT * FROM inbox_items
WHERE source = $1 AND external_ref = $2
LIMIT 1;
```

### 4. Update `CreateInboxItem` query
Extend existing query to accept `external_ref`:
```sql
-- name: CreateInboxItem :one
INSERT INTO inbox_items (content, url, source, item_type, submitted_by, attachments, external_ref)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
```

### 5. Run migration + regen
```bash
# Local test first
psql -h 103.97.125.131 -U postgres -d binhvuongos \
  -f internal/db/migrations/000022_inbox_external_ref.up.sql

# Regen sqlc
sqlc generate
```

## Todo
- [ ] Create `000022_inbox_external_ref.up.sql`
- [ ] Create `000022_inbox_external_ref.down.sql`
- [ ] Append `GetInboxByExternalRef` to `inbox_items.sql`
- [ ] Update `CreateInboxItem` signature (add `external_ref` param)
- [ ] Run `sqlc generate`
- [ ] Verify generated struct in `inbox_items.sql.go` có ExternalRef field
- [ ] Test migration up/down trên dev DB

## Success criteria
- `inbox_items` table có cột `external_ref`
- Unique index enforce no dup `(source, external_ref)`
- Generated Go code có `GetInboxByExternalRef` + `CreateInboxItemParams.ExternalRef`
- Down migration rollback sạch

## Risks
- Existing rows có `external_ref = NULL` → OK (partial index only indexes non-null)
- Production DB migration: chạy manual trên `103.97.125.131` sau khi code merge
