---
type: tester
date: 2026-04-21 00:51 +07
slug: header-size-fix
status: done
---

# QA Report — 431 "Request Header Fields Too Large" + comments loader stuck

## Symptoms reported
- Triage / quick-actions trên inbox trả "Request Header Fields Too Large" (HTTP 431)
- Comments panel kẹt ở "Bình luận: đang tải..."

## Root cause
`fiber.New(fiber.Config{})` không set `ReadBufferSize` → fasthttp default **4096 bytes** cho header buffer.

Browser thực tế gửi: JWT cookie (~280B) + Cloudflare cookies (`cf_clearance` ~1.5KB, `__cf_bm` ~1KB, `_ga`, `_gid`, `cf_analytics`…) + HTMX headers (`HX-Current-URL`, `HX-Target`, `HX-Trigger`, `HX-Request`) + UA/Accept/Referer → dễ vượt 4KB.

Comments panel: `hx-trigger="load"` tự fire GET → request này fail 431 → HTMX không swap → placeholder "Đang tải..." kẹt.

## Fix
`cmd/server/main.go`:
```go
app := fiber.New(fiber.Config{
    AppName:        "Bình Vương OS v2.0",
    ReadBufferSize: 32 * 1024,          // 32KB header buffer
    BodyLimit:      60 * 1024 * 1024,   // 60MB for multipart webhook uploads
})
```

Commit `020f85a` → deploy `103.97.125.186` lúc 00:46 ICT.

## Test results

### A. Endpoint smoke test (15 pages) — ALL PASS
| Path | Status |
|---|---|
| /, /inbox, /dashboard, /work-logs, /tasks, /content, /companies, /campaigns, /knowledge, /bookmarks, /profile, /users, /admin/settings, /notifications, /search | 302/200 ✅ |

### B. Header-size stress test — pre-fix would be 431, now PASS
| Scenario | Result |
|---|---|
| 8KB cookies + HTMX headers → GET /comments | 200 ✅ |
| 8KB cookies + HTMX → POST /inbox/:id/triage | 302 ✅ |
| 16KB cookie payload → GET /inbox | 200 ✅ |

### C. Inbox actions — ALL PASS
| Action | Result |
|---|---|
| GET /comments (HTMX partial) | 200, 752B HTML |
| POST /comments | 200, renders new comment |
| POST /inbox/:id/triage (destination=knowledge) | 302 |
| POST /inbox/:id/archive | 302 |
| GET /inbox/:id (detail) | 200, 13.5KB |

### D. Webhook API (regression check)
| Test | Result |
|---|---|
| `/api/v1/inbox` JSON note | ✅ |
| `/api/v1/inbox` multipart → Drive | ✅ (từ session trước) |
| Idempotent `external_ref` | ✅ |

## Test data cleanup
- Revert `inbox_items[0d05a863…].status` về `raw`, xoá triage_notes
- DELETE comment "QA test comment"

## Unresolved questions

1. Fiber BodyLimit cũ mặc định 4MB — có endpoint nào khác (ngoài webhook multipart) cần upload lớn không? Nếu chỉ webhook thì OK đã set 60MB toàn app.
2. Có muốn add healthcheck endpoint + Uptime monitor (e.g. BetterStack) để catch 5xx trong tương lai sớm không?
3. "HX-Current-URL" rất dài khi chứa query params lồng nhau — cân nhắc giữ 32KB là an toàn trung hạn.
