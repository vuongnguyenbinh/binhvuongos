# Phase 3 — User Detail Admin Page

**Effort:** 45m

## Files

### New
- `web/templates/pages/user_detail.templ`

### Modify
- `internal/handler/users.go` — add `UserDetail` handler
- `cmd/server/main.go` — GET /users/:id
- `web/templates/pages/users.templ` — link row name to detail page

## Handler

```go
func (h *Handler) UserDetail(c *fiber.Ctx) error {
    actor := GetUser(c)
    target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
    if err != nil {
        return c.Status(404).SendString("Không tìm thấy user")
    }
    // Both owner and manager can view; edit gated separately.
    if actor.Role != "owner" && actor.Role != "manager" {
        return c.Status(403).SendString("Forbidden")
    }

    logs, _ := h.queries.ListWorkLogsByUser(c.Context(), target.ID, 10, 0)
    tasks, _ := h.queries.ListTasksByAssignee(c.Context(), target.ID, 10, 0)
    companies, _ := h.queries.ListCompaniesForUser(c.Context(), target.ID)

    return render(c, pages.UserDetailPage(pages.UserDetailData{
        ID:         middleware.UUIDToString(target.ID),
        Email:      target.Email,
        FullName:   target.FullName,
        Role:       target.Role,
        Status:     target.Status,
        Phone:      nullStr(target.Phone),
        AvatarURL:  nullStr(target.AvatarURL),
        CanEdit:    middleware.CanManageUser(actor, target),
        WorkLogs:   toTemplWorkLogs(logs),
        Tasks:      toTemplTaskItems(tasks),
        Companies:  toTemplCompanyOpts(companies),
    }))
}
```

## Queries (may already exist, check)

- `ListWorkLogsByUser(ctx, userID, limit, offset)` — check query/work_logs.sql
- `ListTasksByAssignee(ctx, assigneeID, limit, offset)` — check query/tasks.sql
- `ListCompaniesForUser(ctx, userID)` — via user_company_assignments JOIN companies

Add any missing ones. If conflict with existing helpers, reuse.

## Templ

```templ
templ UserDetailPage(d UserDetailData) {
    @templates.Layout("Chi tiết user — " + d.FullName, "dashboard") {
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- Profile card -->
            <div class="bg-surface border border-hairline rounded-lg p-6">
                if d.AvatarURL != "" {
                    <img src={ d.AvatarURL } class="w-20 h-20 rounded-full mb-3"/>
                }
                <h2 class="display text-xl font-bold">{ d.FullName }</h2>
                <p class="text-sm text-muted mono">{ d.Email }</p>
                <p class="mt-3">Vai trò: <span class="font-semibold">{ roleLabel(d.Role) }</span></p>
                <p>SĐT: { d.Phone }</p>
                <p>Trạng thái: { d.Status }</p>
                if d.CanEdit {
                    <a href={ templ.SafeURL("/users/" + d.ID + "/edit") } class="inline-block mt-4 px-4 py-2 bg-forest text-white rounded text-sm mono">SỬA</a>
                }
            </div>
            <!-- Companies + Tasks + Work logs -->
            <div class="lg:col-span-2 space-y-6">
                <div class="bg-surface border border-hairline rounded-lg p-6">
                    <h3 class="eyebrow mb-3">Công ty phụ trách</h3>
                    ... for _, co := range d.Companies { ... }
                </div>
                <div class="bg-surface border border-hairline rounded-lg p-6">
                    <h3 class="eyebrow mb-3">Công việc gần đây</h3>
                    ... for _, t := range d.Tasks { ... }
                </div>
                <div class="bg-surface border border-hairline rounded-lg p-6">
                    <h3 class="eyebrow mb-3">Work logs gần đây</h3>
                    ... for _, l := range d.WorkLogs { ... }
                </div>
            </div>
        </div>
    }
}
```

## Route

```go
admin.Get("/users/:id", h.UserDetail)  // Note: exists after /users but before /users/:id/edit
```

Fiber precedence: register `/users/:id/edit` first, then `/users/:id` — Fiber matches longer path first? Actually Fiber uses declaration order for same method. Declare `/users/:id` AFTER `/users/:id/edit` + `/users/:id/delete` etc to avoid catching `:id = "edit"`.

## Todo
- [ ] Handler UserDetail
- [ ] Templ UserDetailPage
- [ ] Ensure ListTasksByAssignee / ListWorkLogsByUser / ListCompaniesForUser queries exist (add if missing)
- [ ] Route registered AFTER /users/:id/edit, /users/:id/delete, /users/:id/reset-password
- [ ] Update users.templ row name → link to /users/:id

## Success criteria
- Admin click user row → /users/:id shows profile + tasks + logs + companies
- Staff GET /users/:id → 403
- Manager view owner detail → 200 (just view), but edit button hidden (CanManageUser=false)
