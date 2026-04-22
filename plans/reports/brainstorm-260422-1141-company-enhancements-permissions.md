---
type: brainstorm
date: 2026-04-22 11:41 +07
slug: company-enhancements-permissions
status: approved
---

# Brainstorm: Company Enhancements + Permission Matrix

## Problem statement

1. `/companies` thiếu: edit nâng cao (đã có), **archive/unarchive**, **upload logo**, **filter active/all/archived**, **deadline warning 10 ngày**.
2. Chưa có tài liệu chính thức về phân quyền 3 role (owner/manager/staff) → khó audit + khó mở rộng.

## Decisions (đã chốt)

| Q | A |
|---|---|
| Field deadline | Reuse `companies.end_date` |
| Warning UI | Badge row + Dashboard widget + Notifications (cả 3) |
| Archive scope | Chỉ ẩn company, data liên kết giữ nguyên |
| Ai edit/archive | Owner + manager |
| Notifications schema | Bổ sung `ref_type`/`ref_id` + unique per-day index |
| Dashboard widget visible | Tất cả role; staff filter theo companies assigned |
| Logo formats | image/png, image/jpeg, image/svg+xml |
| Notif goroutine | Startup + ticker 24h |

## Architecture

### Company enhancements

```
/companies?show=active (default) | ?show=all | ?show=archived
  ├─ Filter tabs on top
  ├─ Per-card badge nếu end_date <= today+10d
  ├─ Buttons: [Edit] [Upload Logo] [Archive/Unarchive]
  └─ Query: WHERE ($1='all' OR status=$1)

Archive/Unarchive:
  POST /companies/:id/archive   → status='archived'
  POST /companies/:id/unarchive → status='active'

Logo upload:
  POST /companies/:id/logo (multipart) → drive.UploadFile → logo_url
  Accept: image/png, image/jpeg, image/svg+xml
  Max: 5MB

Deadline badge (helper):
  func companyDeadlineBadge(endDate pgtype.Date) (label string, class string)
    - <0 days: "⚠️ Quá hạn" / rust
    - <=10 days: "⏰ Còn N ngày" / ember
    - else: "" / ""
```

### Dashboard widget

```
[/dashboard page]
  Section "🏢 Công ty sắp hết hạn" (only show if count > 0)
    ├─ Owner/manager: WHERE end_date <= today+10d AND status='active' ORDER BY end_date ASC LIMIT 5
    └─ Staff: JOIN user_company_assignments ON user_id=$actor
```

### Notification cron (in-process goroutine)

```
main.go:
  go startDeadlineNotifier(queries)
    ticker := time.NewTicker(24h)
    runDeadlineCheck() // initial
    for { <-ticker.C; runDeadlineCheck() }

runDeadlineCheck():
  rows = SELECT c.id, c.name, c.end_date, uca.user_id
         FROM companies c JOIN user_company_assignments uca ON uca.company_id=c.id
         WHERE c.status='active' AND c.end_date IS NOT NULL
           AND c.end_date <= CURRENT_DATE + INTERVAL '10 days'
           AND (uca.end_date IS NULL OR uca.end_date > CURRENT_DATE)
  for each (company, user):
    INSERT INTO notifications (user_id, title, body, link, ref_type, ref_id)
    VALUES ($1, 'Công ty sắp hết hạn', ..., '/companies/'||c.id, 'company_deadline', c.id)
    ON CONFLICT (user_id, ref_type, ref_id, notif_date) DO NOTHING;
  +owner always gets notification (no assignment required)
```

### Migration 000024

```sql
-- 000024_notifications_ref.up.sql
ALTER TABLE notifications
  ADD COLUMN ref_type VARCHAR(30),
  ADD COLUMN ref_id   UUID,
  ADD COLUMN notif_date DATE NOT NULL DEFAULT CURRENT_DATE;

CREATE UNIQUE INDEX idx_notifications_dedup_per_day
  ON notifications(user_id, ref_type, ref_id, notif_date)
  WHERE ref_type IS NOT NULL;
```

### Permission matrix doc

File: `docs/permissions.md`

Format: bảng markdown với cột `Owner | Manager | Staff`, hàng = action per resource. Kèm:
- Legend (✅ / ⚠️ conditional / ❌)
- Notes về cách enforce (`RequireRole`, `CanManageUser`)
- File paths reference cho mỗi enforcement point

## Effort

| # | Phase | LOC | Time |
|---|---|---|---|
| 1 | Migration 000024 + Notification model updates + CreateNotification dedup | 60 | 45m |
| 2 | Archive/unarchive handlers + route + queries (UpdateCompanyStatus) | 50 | 30m |
| 3 | Filter UI (3-tab segmented) + query variants (ListCompanies by status) | 70 | 30m |
| 4 | Logo upload handler + company card/detail update | 70 | 45m |
| 5 | Deadline badge helper + render on list + detail | 50 | 30m |
| 6 | Dashboard widget (owner/manager all, staff filtered) | 80 | 45m |
| 7 | Deadline notifier goroutine in main.go | 70 | 45m |
| 8 | docs/permissions.md matrix | markdown | 45m |
| 9 | Deploy + E2E test | — | 30m |
| **Total** | | **~450 LOC** | **~5h 15m** |

## Risks

| Risk | Mitigation |
|---|---|
| Goroutine panic crashes main | Wrap `runDeadlineCheck()` in `recover()` |
| Large company list causes spike at startup | Batch process; current scale ≤10 companies → trivial |
| Notif dedup race on multiple replicas | Unique index at DB layer prevents dup regardless of app replicas |
| Staff sees company they're NOT assigned to | JOIN filter enforces; owner bypass explicitly |
| Logo SVG with embedded script (XSS) | Render via `<img>` tag only (no `<object>` / `<iframe>`) — SVG rendered as image is safe |
| Timezone-sensitive "còn N ngày" wrong by 1 | Use `CURRENT_DATE` (DB server TZ) + Vietnam TZ set in compose/env |

## Success criteria

- [ ] `/companies?show=archived` hiển thị đúng 0 row ban đầu
- [ ] Owner click Archive → row ẩn khỏi default list, có trong `?show=archived`
- [ ] Manager upload logo PNG → thấy ngay trên card
- [ ] SVG/JPG upload OK; GIF/PDF từ chối
- [ ] Tạo 1 test company với `end_date = today+5` → badge "⏰ Còn 5 ngày"
- [ ] Dashboard widget hiển thị companies gần deadline
- [ ] Notifications table có row `ref_type='company_deadline'` sau khi goroutine chạy
- [ ] 2 lần chạy goroutine cùng ngày → 1 row (unique)
- [ ] Staff chỉ thấy widget companies họ được assign
- [ ] Webhook + user CRUD + inbox regression OK

## Unresolved

Không còn câu hỏi. Design complete.
