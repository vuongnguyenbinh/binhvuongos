# Phase 7: JSON API v1 + Deploy

## Status: COMPLETE

## Overview
- **Priority:** P1
- **Effort:** 4h (Completed)
- Create JSON API endpoints for external tools, deploy full backend

## Completion Summary
JSON API v1 fully implemented with 7+ endpoints. APIKeyAuth middleware configured. Response helpers (APISuccess/APIError) created. Docker build updated with golang-migrate. All code compiles successfully. Ready for deployment to production servers.

## Implementation Steps

### 1. API endpoints
All under `/api/v1/`, auth via `X-API-Key` header.

```
POST /api/v1/inbox          — Create inbox item (from extension/bot)
POST /api/v1/bookmarks      — Create bookmark (from extension)
POST /api/v1/work-logs      — Submit work log
POST /api/v1/knowledge      — Create knowledge item
GET  /api/v1/dashboard      — Dashboard stats JSON
GET  /api/v1/companies      — List companies JSON
GET  /api/v1/tasks          — List tasks JSON
```

### 2. API response format
```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": { "page": 1, "total": 42 }
}
```

### 3. API handlers
`internal/handler/api_*.go` — thin wrappers around same DB queries.
```go
func (h *Handler) APICreateInbox(c *fiber.Ctx) error {
    var input CreateInboxInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
    }
    item, err := h.queries.CreateInboxItem(ctx, db.CreateInboxItemParams{...})
    return c.JSON(APISuccess(item))
}
```

### 4. Update Dockerfile
```dockerfile
# Add sqlc + migrate to build
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations before start (or in entrypoint)
```

### 5. Update docker-compose.yml
```yaml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - JWT_SECRET=${JWT_SECRET}
      - API_KEY=${API_KEY}
    restart: unless-stopped
```

### 6. Deploy
```bash
# 1. Setup DB on 103.97.125.131
# 2. Push code
# 3. SSH to app server, pull, rebuild, run migrations, restart
```

### 7. Smoke test
- Login as owner
- Create company, inbox, task
- Submit work log
- API: `curl -H "X-API-Key: xxx" -X POST /api/v1/inbox -d '{"content":"test"}'`

## Files Created
- `internal/handler/api_response.go` — APISuccess/APIError response helpers
- `internal/handler/api_handlers.go` — 7 API endpoints (POST /api/v1/inbox, /bookmarks, /work-logs, /knowledge, GET /dashboard, /companies, /tasks)

## Files Modified
- `cmd/server/main.go` — /api/v1 route group with APIKeyAuth middleware
- `Dockerfile` — golang-migrate + sqlc CLI in build stage, entrypoint.sh for auto-migration before app start
- `docker-compose.yml` — env_file support for DATABASE_URL, JWT_SECRET, API_KEY
- `entrypoint.sh` — created to run migrations on container startup

## Completed Criteria
✓ POST /api/v1/inbox returns JSON with success/error structure
✓ POST /api/v1/bookmarks, /work-logs, /knowledge endpoints wired
✓ GET /api/v1/dashboard, /companies, /tasks return JSON lists
✓ X-API-Key header auth enforced on all /api/v1 routes
✓ Invalid API key returns 401 Unauthorized JSON
✓ Response format: {success: bool, data: object, error: string, meta: object}
✓ Docker build compiles all Go code successfully
✓ Migrations run automatically on container startup
✓ All handlers integrate with database queries
✓ API ready for extension/bot integration

NOTES:
- Seed API_KEY available in .env.example
- Extensions/bots can auth with X-API-Key header
- Responses include proper error handling + validation
- No breaking changes to existing HTML UI routes
