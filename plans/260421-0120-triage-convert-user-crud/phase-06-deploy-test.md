# Phase 6 — Deploy + E2E Tests

**Effort:** 30m | **Priority:** P1 | **Depends on:** Phase 1-5

## Pre-deploy

1. `go build ./...` + `go vet ./...` pass locally (via Docker golang:1.22-alpine)
2. Dry-run migration 000023 on dev psql:
   ```bash
   ssh root@103.97.125.131 \
     'PGPASSWORD=bvos_secure_2026 psql -U bvos -d binhvuong_os -h localhost -c "SELECT role, COUNT(*) FROM users GROUP BY role;"'
   ```
   Expect: `owner | 1`, `core_staff | 2`, `ctv | 2` → after migration: `owner | 1`, `manager | 2`, `staff | 2`
3. Backup users table:
   ```bash
   pg_dump -t users > /tmp/users-backup-pre-000023.sql
   ```

## Deploy

```bash
git add -A
git commit -m "feat: triage convert + user CRUD + 3-role consolidation"
git push

sshpass -p 'pi8phohVieyai9i' ssh root@103.97.125.186 \
  "cd /opt/binhvuongos && git pull && \
   docker rmi binhvuongos-app 2>/dev/null; \
   docker build --no-cache -t binhvuongos-app . && \
   docker compose up -d"
```

Entrypoint tự chạy migration 000023. Verify log:
```
Running database migrations...
23/u role_consolidation (...)
```

## E2E tests

### A. Role migration
```bash
ssh root@103.97.125.131 \
  'PGPASSWORD=bvos_secure_2026 psql -U bvos -d binhvuong_os -h localhost -c "SELECT role, COUNT(*) FROM users GROUP BY role;"'
# Expect: only owner, manager, staff
```

### B. Triage → Task (as owner)
```bash
# Login owner, get cookie
curl -c /tmp/c.txt -d 'email=vuongnguyenbinh@gmail.com&password=BinhVuong2026!' \
  https://os.binhvuong.vn/auth/login

# Find raw inbox item ID
IID=$(psql ... "SELECT id FROM inbox_items WHERE status='raw' LIMIT 1")

# Convert to task
curl -b /tmp/c.txt -X POST "https://os.binhvuong.vn/inbox/$IID/convert?target=task" \
  -d "title=E2E%20task&priority=normal&triage_notes=e2e"

# Verify
psql ... "SELECT id, title FROM tasks WHERE title='E2E task';"
psql ... "SELECT status, converted_to_type, converted_to_id FROM inbox_items WHERE id='$IID';"
# Expect: status=done, converted_to_type=task, converted_to_id=<task.id>
```

### C. Triage → Content (requires company_id)
- POST without company_id → 400 error
- POST with valid company_id → 201 + inbox done

### D. Triage → Knowledge
- POST with title + body → success; row in knowledge_items

### E. Idempotent convert (double-submit)
- Same inbox item, 2 parallel POST → 1 succeeds, 2nd returns 409

### F. User CRUD (as manager)
1. Login as manager user
2. `POST /users` with role=staff → 201
3. `POST /users` with role=manager → 403 (whitelist reject)
4. `POST /users/<staff_id>` update → 200
5. `POST /users/<manager_id>` update (edit peer manager) → 403
6. `POST /users/<staff_id>/delete` → 200; verify `status=archived`

### G. User CRUD (as owner)
1. Create manager → 201
2. Create staff → 201
3. Edit any user → 200
4. Delete self → 400 (self-protect)

### H. Password reset
1. Owner click Reset on staff → URL returned in flash
2. GET `/reset/<token>` → form render
3. POST new password (<8 chars) → 400
4. POST valid password → 302 /login
5. Login with new password → success
6. Reuse same token → 404 used
7. Token >1h old → 404 expired (manually set expires_at)

### I. Auth + status
1. Soft-delete a user
2. Their cached session → GET any page → 302 /login (active check)

### J. No regression — webhook still works
```bash
curl -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" -H 'Content-Type: application/json' \
  -d '{"content":"post-deploy test","source":"manual"}'
# Expect 201
```

## Cleanup

```bash
# Remove test data
psql ... "DELETE FROM tasks WHERE title='E2E task';"
psql ... "DELETE FROM content WHERE title LIKE 'E2E%';"
psql ... "DELETE FROM knowledge_items WHERE title LIKE 'E2E%';"
# Revert inbox items used in tests back to status=raw if needed
```

## Rollback plan

If something breaks post-deploy:
```bash
# Revert code
ssh root@103.97.125.186 \
  "cd /opt/binhvuongos && git reset --hard HEAD~N && docker compose restart app"

# Rollback migration 000023
ssh root@103.97.125.131 \
  'PGPASSWORD=bvos_secure_2026 psql -U bvos -d binhvuong_os -h localhost \
   -f /tmp/000023_role_consolidation.down.sql'

# Restore user data from backup if oversweep happened
psql ... < /tmp/users-backup-pre-000023.sql
```

## Todo
- [ ] Pre-deploy dry-run migration check
- [ ] Backup users table
- [ ] git commit + push
- [ ] Server deploy (docker rebuild + up)
- [ ] Verify migration log `23/u`
- [ ] Run E2E A-J suite, document failures
- [ ] Cleanup test data
- [ ] Update CLAUDE.md if env/deploy procedure changed
- [ ] Mark plan.md `status=completed`

## Success criteria
- All 10 E2E groups pass
- Zero regression on webhook + inbox list + comments
- Migration verified via role count query
- No leftover test data in DB

## Risks
- Docker build >3min on server; if build fail user-facing 503 briefly. Mitigation: build before down (already default compose behavior via `--no-cache` then `up -d`)
- Session mismatch window: users with JWT role=core_staff hit server with new code → middleware rejects. Mitigation: all active users re-login. Expected <5 users, acceptable.
