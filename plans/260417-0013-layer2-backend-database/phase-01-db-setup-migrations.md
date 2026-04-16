# Phase 1: PostgreSQL Setup + Migrations

## Status: COMPLETE

## Overview
- **Priority:** P0 — blocks everything
- **Effort:** 4h (Completed)
- Setup PostgreSQL database on 103.97.125.131, create 12 migration files from existing schema

## Completion Summary
All database migrations implemented and tested. 15 migration pairs created covering all 12 tables plus bookmarks, views, and seed data. Docker build includes golang-migrate with auto-migration on startup.

## Implementation Steps

### 1. SSH to DB server, create database
```bash
ssh root@103.97.125.131
sudo -u postgres psql
CREATE DATABASE binhvuong_os;
CREATE USER bvos WITH PASSWORD 'secure_password_here';
GRANT ALL PRIVILEGES ON DATABASE binhvuong_os TO bvos;
\c binhvuong_os
GRANT ALL ON SCHEMA public TO bvos;
```

### 2. Verify network connectivity from app server
```bash
# From 103.97.125.186
psql -h 103.97.125.131 -U bvos -d binhvuong_os
```
If blocked: edit `pg_hba.conf` on DB server to allow `103.97.125.186`.

### 3. Create .env file
```env
DATABASE_URL=postgres://bvos:password@103.97.125.131:5432/binhvuong_os?sslmode=disable
JWT_SECRET=random_64_char_string
API_KEY=bvos_api_random_32_chars
PORT=3000
```

### 4. Create config module
`internal/config/config.go` — load env vars with godotenv.

### 5. Create migration files
From `Root/binhvuong-schema-postgres.md`, extract into sequential migration files:

| # | File | Tables |
|---|------|--------|
| 1 | 000001_extensions | pgcrypto, pg_trgm, unaccent + trigger function |
| 2 | 000002_users | users + indexes + trigger |
| 3 | 000003_companies | companies + indexes + trigger |
| 4 | 000004_user_company_assignments | assignments + indexes + trigger |
| 5 | 000005_inbox_items | inbox_items + indexes + trigger |
| 6 | 000006_campaigns | campaigns + indexes + trigger |
| 7 | 000007_tasks | tasks + indexes + trigger (FK to campaigns) |
| 8 | 000008_work_types | work_types + seed data + trigger |
| 9 | 000009_content | content + indexes + FK tasks.content_id + trigger |
| 10 | 000010_work_logs | work_logs + views + indexes + trigger |
| 11 | 000011_knowledge_items | knowledge_items + indexes + FTS + trigger |
| 12 | 000012_objectives | objectives + indexes + trigger |
| 13 | 000013_notion_sync_log | notion_sync_log + indexes |
| 14 | 000014_seed_owner | INSERT owner user (bcrypt hashed password) |

Each migration has `.up.sql` and `.down.sql`.

### 6. Run migrations
```bash
migrate -database "${DATABASE_URL}" -path internal/db/migrations up
```

### 7. Update Dockerfile
Add golang-migrate + sqlc to Docker build. Add .env support.

### 8. Update docker-compose.yml
Add environment variables for DATABASE_URL.

## Files Created
- `.env` + `.env.example`
- `internal/config/config.go` — loads DB_URL, JWT_SECRET, API_KEY, PORT from .env
- `internal/db/migrations/000001-000015` (30 SQL files)
  - 000001: extensions (pgcrypto, pg_trgm, unaccent)
  - 000002-000012: all 12 tables (users, companies, assignments, inbox, tasks, content, campaigns, work_types, work_logs, knowledge, objectives, notion_sync_log)
  - 000013: bookmarks table
  - 000014: seed owner user (owner@binhvuong.vn / BinhVuong2026!)
  - 000015: views for work_logs daily/monthly stats and campaign progress

## Files Modified
- `go.mod` — added pgx, godotenv, jwt (golang-jwt/jwt/v5), bcrypt
- `Dockerfile` — golang-migrate + sqlc CLI in build stage, entrypoint.sh for auto-migration
- `docker-compose.yml` — env_file support
- `.gitignore` — added .env

## Completed Criteria
✓ `migrate up` creates all 12 tables + views
✓ Owner user seeded with bcrypt-hashed password
✓ App connects to remote DB from Docker container
✓ Migration files include indexes + triggers on updated_at columns
✓ Foreign keys + cascading deletes configured
