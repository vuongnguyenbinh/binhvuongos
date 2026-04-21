# Phase 5 — HTMX Triage Modals + UI Wiring

**Effort:** 75m | **Priority:** P1 | **Depends on:** Phase 4

## Context
- Existing component: `web/templates/components/modal.templ` (check pattern)
- Existing quick-action buttons in `web/templates/pages/inbox.templ:150-160` — use `<form method="POST">` directly
- Need: replace with HTMX modal triggers

## Files

### Create
- `web/templates/components/triage_modals.templ` — 3 modals: `TriageTaskModal`, `TriageContentModal`, `TriageKnowledgeModal`

### Modify
- `web/templates/pages/inbox.templ` — replace submit forms with HTMX triggers `hx-get` to load modal
- `web/templates/pages/inbox_detail.templ` — same for detail page triage panel
- `cmd/server/main.go` — add `GET /inbox/:id/triage-modal?target=xxx` to render modal HTML
- `internal/handler/inbox_convert.go` — add `TriageModalPartial(c)` returns HTMX partial

## UI pattern

```templ
// inbox.templ row action
<button
    hx-get={ fmt.Sprintf("/inbox/%s/triage-modal?target=task", item.ID) }
    hx-target="#modal-root"
    hx-swap="innerHTML"
    class="...">
    ✅
</button>

// Global modal container in layout.templ:
<div id="modal-root"></div>
```

## Modal template (example: Task)

```templ
templ TriageTaskModal(inboxID string, contentPrefill string, companies []CompanyOpt) {
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" hx-on:click="this.remove()">
        <div class="bg-surface rounded-lg p-6 w-full max-w-md" hx-on:click="event.stopPropagation()">
            <h3 class="display text-lg font-bold mb-4">→ Tạo công việc</h3>
            <form hx-post={ "/inbox/" + inboxID + "/convert?target=task" }
                  hx-headers='{"Accept":"text/html"}'
                  hx-target="body"
                  hx-swap="none"
                  hx-on::after-request="if(event.detail.successful) window.location='/inbox'">
                <label class="block text-sm mb-1">Tiêu đề</label>
                <input type="text" name="title" required
                    value={ truncate(contentPrefill, 80) }
                    class="w-full px-3 py-2 border rounded mb-3"/>

                <label class="block text-sm mb-1">Công ty</label>
                <select name="company_id" class="w-full px-3 py-2 border rounded mb-3">
                    <option value="">— Không chọn —</option>
                    for _, c := range companies {
                        <option value={ c.ID }>{ c.Name }</option>
                    }
                </select>

                <label class="block text-sm mb-1">Ưu tiên</label>
                <select name="priority" class="w-full px-3 py-2 border rounded mb-3">
                    <option value="normal">Bình thường</option>
                    <option value="high">Cao</option>
                    <option value="low">Thấp</option>
                </select>

                <label class="block text-sm mb-1">Hạn</label>
                <input type="date" name="due_date" class="w-full px-3 py-2 border rounded mb-4"/>

                <label class="block text-sm mb-1">Ghi chú triage</label>
                <textarea name="triage_notes" rows="2" class="w-full px-3 py-2 border rounded mb-4"></textarea>

                <div class="flex gap-2 justify-end">
                    <button type="button" hx-on:click="document.getElementById('modal-root').innerHTML=''"
                        class="px-4 py-2 text-muted">Hủy</button>
                    <button type="submit" class="px-4 py-2 bg-forest text-white rounded">Tạo & đánh dấu xong</button>
                </div>
            </form>
        </div>
    </div>
}
```

`TriageContentModal`: title, content_type (blog/social/video), company_id (**REQUIRED**), triage_notes.
`TriageKnowledgeModal`: title, body (textarea, prefill content), category, triage_notes.

## Handler: modal loader

```go
func (h *Handler) TriageModalPartial(c *fiber.Ctx) error {
    target := c.Query("target")
    inboxID := c.Params("id")
    item, err := h.queries.GetInboxItemByID(c.Context(), middleware.StringToUUID(inboxID))
    if err != nil { return c.Status(404).SendString("Không tìm thấy") }

    companies, _ := h.queries.ListCompanies(c.Context(), 100, 0)
    opts := toCompanyOpts(companies)

    switch target {
    case "task":      return render(c, components.TriageTaskModal(inboxID, item.Content, opts))
    case "content":   return render(c, components.TriageContentModal(inboxID, item.Content, opts))
    case "knowledge": return render(c, components.TriageKnowledgeModal(inboxID, item.Content))
    default:          return c.Status(400).SendString("Target không hợp lệ")
    }
}
```

## Todo
- [ ] Create `components/triage_modals.templ` with 3 variants
- [ ] Add `#modal-root` div in `layout.templ`
- [ ] Replace form-submit quick actions in `inbox.templ` + `inbox_detail.templ` with HTMX modal triggers
- [ ] Add route `GET /inbox/:id/triage-modal`
- [ ] Add handler `TriageModalPartial`
- [ ] Handle content modal: `company_id` required (client + server validate)
- [ ] `templ generate` + `go build ./...` pass

## Success criteria
- Click `→ Task` → modal xuất hiện với prefill title
- Submit modal thành công → redirect `/inbox`, inbox item biến mất khỏi list
- Click ngoài modal hoặc Hủy → modal đóng
- Content modal không chọn company → 400 với flash error
- ESC key hoặc overlay click đóng modal

## Risks
- HTMX + templ render cần `templ generate` sau mọi thay đổi; Docker build tự chạy nhưng local dev phải nhớ
- Modal stacking nếu user click nhanh 2 target khác nhau → solve by replacing `#modal-root innerHTML` (last wins)
- Mobile: modal max-width 100vw, padding cần responsive
