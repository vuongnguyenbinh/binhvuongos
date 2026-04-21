# Phase 2 — User CRUD + Permission Middleware

**Effort:** 90m | **Priority:** P1 | **Depends on:** Phase 1

## Context
- Existing: `handler/users.go` has `Users` (list) + `CreateUser` (insecure — trusts form role)
- Missing: Update, Delete, dedicated edit page, permission enforcement

## Files

### Create
- `internal/middleware/user_perm.go` — `CanManageUser(actor, target) bool` + helper `AllowedTargetRoles(actor) []string`
- `internal/handler/user_crud.go` — `UpdateUser`, `DeleteUser` (soft), `EditUserPage`
- `web/templates/pages/user_edit.templ` — edit form page

### Modify
- `internal/handler/users.go::CreateUser` — enforce role whitelist per actor
- `internal/db/query/users.sql` + `generated/users.sql.go` — add `UpdateUser`, `SoftDeleteUser`, `GetUserByID`
- `web/templates/pages/users.templ` — role-gated edit/delete buttons
- `cmd/server/main.go` — register new routes

## Permission matrix

```go
// middleware/user_perm.go
func CanManageUser(actor generated.User, target generated.User) bool {
    if actor.Role == "owner" { return true }
    if actor.Role == "manager" && target.Role == "staff" { return true }
    return false
}

func AllowedTargetRoles(actor generated.User) []string {
    switch actor.Role {
    case "owner":   return []string{"owner", "manager", "staff"}
    case "manager": return []string{"staff"}
    default:        return nil
    }
}
```

## Routes

```go
admin := app.Group("", middleware.RequireRole("owner", "manager"))
admin.Get("/users",          h.Users)
admin.Post("/users",         h.CreateUser)          // enforce whitelist
admin.Get("/users/:id/edit", h.EditUserPage)
admin.Post("/users/:id",     h.UpdateUser)
admin.Post("/users/:id/delete", h.DeleteUser)
```

## Queries

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET full_name=$2, role=$3, phone=$4, status=$5
WHERE id=$1 RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users SET status='archived' WHERE id=$1;
```

## Handler contract

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    actor := GetUser(c)
    targetRole := c.FormValue("role")
    if !slices.Contains(middleware.AllowedTargetRoles(actor), targetRole) {
        return c.Status(403).SendString("Không có quyền tạo user với role này")
    }
    // ... rest
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
    actor := GetUser(c)
    target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
    if err != nil { return c.Status(404).SendString("Không tìm thấy") }
    if !middleware.CanManageUser(actor, target) {
        return c.Status(403).SendString("Không có quyền")
    }
    newRole := c.FormValue("role")
    if !slices.Contains(middleware.AllowedTargetRoles(actor), newRole) {
        return c.Status(403).SendString("Không có quyền set role này")
    }
    // ... UPDATE
}

// DeleteUser — similar perm check, plus guard against self-delete
```

## Auth middleware update (status check)

```go
// middleware/auth.go — add after loading user:
if user.Status != "active" {
    c.ClearCookie("token")
    return c.Redirect("/login")
}
```

## Todo
- [ ] Create `middleware/user_perm.go`
- [ ] Add `GetUserByID`, `UpdateUser`, `SoftDeleteUser` in `generated/users.sql.go` + `query/users.sql`
- [ ] Create `handler/user_crud.go` with 3 handlers
- [ ] Create `web/templates/pages/user_edit.templ`
- [ ] Update `CreateUser` with whitelist enforcement + self-protection
- [ ] Add routes in `main.go`
- [ ] Role-gate edit/delete buttons in `users.templ`
- [ ] Update `AuthRequired` to reject non-active users
- [ ] `go build ./...` pass

## Success criteria
- Manager POST /users with role=manager → 403
- Manager POST /users/:owner_id → 403
- Owner POST /users with role=manager → 201 + redirect
- Deleted user → `status=archived`, next request 302 login
- Self-delete owner → blocked 400

## Risks
- Forgot to check `target.Status` in edit page → allow editing deleted users (accept)
- Role whitelist check using `slices.Contains` — Go 1.21+. Verify `go.mod` (project on 1.22 OK)
