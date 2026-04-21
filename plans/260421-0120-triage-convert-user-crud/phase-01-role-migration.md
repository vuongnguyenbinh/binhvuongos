# Phase 1 — Role Migration + i18n

**Effort:** 30m | **Priority:** P1 | **Status:** pending

## Context
- Current roles in DB: `owner` (1), `core_staff` (2), `ctv` (2), `staff` (default)
- Role enforcement: `internal/middleware/role.go::RequireRole(...string)`
- Vietnamese labels: `internal/handler/i18n.go` (check existing)

## Files

### Create
- `internal/db/migrations/000023_role_consolidation.up.sql`
- `internal/db/migrations/000023_role_consolidation.down.sql`

### Modify
- `internal/handler/i18n.go` — add/update `roleLabel(role string) string`
- `cmd/server/main.go` — update `RequireRole("owner","core_staff")` → `RequireRole("owner","manager")` everywhere
- `web/templates/pages/users.templ` — role dropdown options
- `web/templates/pages/dashboard.templ` — if role comparison hardcoded

## Steps

### 1. Write migration
```sql
-- 000023_role_consolidation.up.sql
-- Remap: core_staff → manager, ctv → staff
UPDATE users SET role='manager' WHERE role='core_staff';
UPDATE users SET role='staff'   WHERE role IN ('ctv');
-- any other legacy values default to staff
UPDATE users SET role='staff'
  WHERE role NOT IN ('owner', 'manager', 'staff');

ALTER TABLE users ADD CONSTRAINT chk_user_role
  CHECK (role IN ('owner', 'manager', 'staff'));
```

```sql
-- 000023_role_consolidation.down.sql
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_user_role;
-- No data rollback — forward-only migration
```

### 2. Update `RequireRole` call sites
```bash
grep -rn 'RequireRole(' cmd/ internal/ | grep -v _test.go
```
Replace `"core_staff"` → `"manager"` everywhere. Example:
```go
// main.go
admin := app.Group("", middleware.RequireRole("owner", "manager"))
```

### 3. i18n role labels
```go
// internal/handler/i18n.go
func RoleLabel(role string) string {
    switch role {
    case "owner":   return "Chủ sở hữu"
    case "manager": return "Quản lý"
    case "staff":   return "Nhân sự"
    default:        return role
    }
}
```

### 4. UI dropdown
```templ
// users.templ — create form role select
<select name="role">
    <option value="staff">Nhân sự</option>
    if isOwner {
        <option value="manager">Quản lý</option>
    }
</select>
```

## Todo
- [ ] Write `000023_role_consolidation.up.sql`
- [ ] Write `000023_role_consolidation.down.sql`
- [ ] Update `RequireRole` call sites (grep+sed safe)
- [ ] Add `RoleLabel` helper in `i18n.go`
- [ ] Replace hardcoded role strings in templates with `RoleLabel()`
- [ ] Dry-run migration on dev DB, verify counts
- [ ] `go build ./...` pass

## Success criteria
- `SELECT role, COUNT(*) FROM users GROUP BY role;` returns only `owner`, `manager`, `staff`
- `RequireRole("owner","manager")` blocks `staff` users from `/users`, `/admin/*`
- UI shows Vietnamese labels everywhere

## Risks
- Migration chạy trên prod khi có user đang login với old role string → JWT vẫn có old role → middleware fail. Mitigation: force logout all sessions (rotate JWT secret) HOẶC accept short window mismatch.
