# Phase 3 — Cleanup Old Telegram Handler

**Effort:** 20m | **Priority:** P2 | **Status:** pending | **Depends on:** Phase 2

## Context
- Old handler: `internal/handler/integrations.go::TelegramWebhook` (lines 34-91)
- Old route: `cmd/server/main.go:58` — `api.Post("/telegram/webhook", h.TelegramWebhook)`
- Unused query: `queries.GetUserByTelegramID` (only called from `TelegramWebhook`)

## Overview
Xoá handler Telegram cũ (broken do `APIKeyAuth` wrap). Telegram bot sẽ đi qua n8n → unified `/api/v1/inbox`. Giữ `users.telegram_id` column (có thể vẫn dùng để lookup user trong n8n flow).

## Files to modify
- `internal/handler/integrations.go` — delete `TelegramWebhook`, `isURL` (nếu không dùng nơi khác)
- `cmd/server/main.go` — delete route line

## Files to keep (không xoá)
- `users.telegram_id` column — n8n có thể tra cứu
- `queries.GetUserByTelegramID` — có thể dùng trong future n8n function, giữ

## Implementation steps

### 1. Check `isURL` usage elsewhere
```bash
grep -rn "isURL\|isUrl" internal/ cmd/
```
Nếu chỉ dùng trong `TelegramWebhook` → xoá cùng. Nếu có nơi khác → giữ, chỉ xoá `TelegramWebhook`.

### 2. Remove `TelegramWebhook` handler
Xoá function `TelegramWebhook` (integrations.go:34-91).

### 3. Remove route
`cmd/server/main.go`: xoá dòng `api.Post("/telegram/webhook", h.TelegramWebhook)`.

### 4. Remove unused imports trong integrations.go
- `database/sql` có thể bỏ nếu không còn dùng
- `fmt`, `log` tương tự

### 5. Compile check
```bash
go build ./...
```

## Todo
- [ ] Grep `isURL` trong codebase — quyết định xoá hay giữ
- [ ] Delete `TelegramWebhook` function
- [ ] Delete route registration in `main.go`
- [ ] Clean unused imports in `integrations.go`
- [ ] `go build ./...` pass
- [ ] Test: `curl -X POST https://os.binhvuong.vn/api/v1/telegram/webhook` → 404

## Success criteria
- `/api/v1/telegram/webhook` trả 404 (route bị xoá)
- `integrations.go` chỉ còn `NotionSyncStatus` + `NotionSyncTrigger` (Notion stubs giữ lại)
- `go build ./...` pass
- Không còn reference tới `TelegramWebhook` trong codebase

## Risks
- n8n flow Telegram chưa setup xong ở phase 5 → tạm thời không có incoming từ Telegram. **Acceptable** vì handler cũ đã broken từ đầu.
