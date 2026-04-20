# Phase 5 — Deploy + E2E Test

**Effort:** 30m | **Priority:** P1 | **Status:** pending | **Depends on:** Phase 1-4

## Context
- Server: `103.97.125.186` (root / pi8phohVieyai9i)
- DB: `103.97.125.131` (PostgreSQL)
- Deploy script pattern từ CLAUDE.md:
  ```bash
  ssh root@103.97.125.186 "cd /opt/binhvuongos && git pull && \
    docker rmi binhvuongos-app 2>/dev/null; \
    docker build --no-cache -t binhvuongos-app . && \
    docker compose up -d"
  ```

## Overview
Generate master API key, update server `.env`, run migration, deploy code, test end-to-end với cURL + n8n sample.

## Implementation steps

### 1. Generate master API key (local)
```bash
API_KEY=$(openssl rand -hex 32)
echo "Generated key: $API_KEY"
# SAVE to password manager
```

### 2. Update server .env
```bash
ssh root@103.97.125.186 << EOF
cd /opt/binhvuongos
# Append or update
sed -i '/^API_KEY=/d' .env
echo "API_KEY=$API_KEY" >> .env
# OWNER_EMAIL nếu chưa có
grep -q OWNER_EMAIL .env || echo "OWNER_EMAIL=vuongnguyenbinh@gmail.com" >> .env
cat .env | grep -E "API_KEY|OWNER_EMAIL"
EOF
```

### 3. Run migration trên production DB
```bash
# Từ server app, chạy migrate tool
ssh root@103.97.125.186 << 'EOF'
cd /opt/binhvuongos
docker compose exec app migrate -path /app/migrations -database "$DATABASE_URL" up
EOF
```
(Nếu app image không có migrate binary — chạy psql trực tiếp)

### 4. Deploy code
```bash
# Từ local
git add -A
git commit -m "feat: unified inbox webhook API with multipart + idempotency"
git push

# Trên server
sshpass -p 'pi8phohVieyai9i' ssh root@103.97.125.186 \
  "cd /opt/binhvuongos && git pull && \
   docker rmi binhvuongos-app 2>/dev/null; \
   docker build --no-cache -t binhvuongos-app . && \
   docker compose up -d"
```

Bump cache-busting `?v=` nếu có sửa CSS/JS (không có trong plan này, skip).

### 5. E2E test scripts

#### a. Health check
```bash
curl -s https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: wrong" -X POST -d '{}' \
  | jq  # expect 401
```

#### b. JSON note
```bash
curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"content":"Test note từ cURL","source":"manual"}' | jq
# expect 201, item_type=note
```

#### c. JSON link (auto-detect)
```bash
curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"content":"https://ahrefs.com/blog","source":"manual"}' | jq
# expect item_type=link, url auto-populated
```

#### d. Multipart with file
```bash
echo "hello" > /tmp/test.txt
curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" \
  -F "content=File test" -F "source=manual" -F "item_type=file" \
  -F "file=@/tmp/test.txt" | jq
# expect attachments[0].url = drive.google.com/...
```

#### e. Idempotency
```bash
# Send twice with same external_ref
for i in 1 2; do
  curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
    -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
    -d '{"content":"idempotent test","source":"manual","external_ref":"test-001"}' | jq '.duplicate'
done
# expect: null then true
```

#### f. Validation errors
```bash
# Empty content
curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
  -d '{"source":"manual"}' | jq
# expect 400

# Invalid item_type
curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
  -d '{"content":"x","source":"manual","item_type":"invalid"}' | jq
# expect 400

# 11 attachments
curl -s -X POST https://os.binhvuong.vn/api/v1/inbox \
  -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
  -d '{"content":"x","source":"manual","attachment_urls":["a","b","c","d","e","f","g","h","i","j","k"]}' | jq
# expect 400
```

#### g. Old Telegram route removed
```bash
curl -s -o /dev/null -w "%{http_code}\n" \
  -X POST https://os.binhvuong.vn/api/v1/telegram/webhook
# expect 404
```

### 6. Verify DB
```bash
ssh root@103.97.125.131 'psql -U postgres -d binhvuongos -c "SELECT id, source, item_type, external_ref, jsonb_array_length(attachments) as n_att, created_at FROM inbox_items ORDER BY created_at DESC LIMIT 10;"'
```

### 7. Sample n8n flow smoke test
- Import `docs/n8n-flows/telegram-to-inbox.json`
- Set `BVOS_API_KEY` env var in n8n
- Set Telegram bot credentials
- Activate workflow
- Send 3 messages to bot from Telegram
- Verify 3 rows with `source=telegram` appear (no dup)

## Todo
- [ ] Generate API_KEY, save password manager
- [ ] Update `.env` trên server
- [ ] Run migration `000022` trên production DB
- [ ] git commit + push + deploy script
- [ ] Verify startup log: "Owner user resolved: <uuid>"
- [ ] Run E2E tests a-g, all pass
- [ ] Verify DB rows qua psql
- [ ] Setup n8n Telegram flow (if n8n available), smoke test 3 msg
- [ ] Verify old `/api/v1/telegram/webhook` returns 404
- [ ] Update `CLAUDE.md` nếu có thêm env vars mới

## Success criteria
- Tất cả 7 E2E tests pass
- DB có rows với đầy đủ fields (source, item_type, external_ref, attachments, submitted_by)
- n8n flow gửi 3 msg → 3 rows (idempotent, no dup)
- `/api/v1/telegram/webhook` 404
- Docs readable, người ngoài setup được 1 source mới

## Rollback plan
Nếu deploy fail:
```bash
# Revert commit
ssh root@103.97.125.186 "cd /opt/binhvuongos && git reset --hard HEAD~1 && docker compose restart app"

# Rollback migration
ssh root@103.97.125.131 'psql -U postgres -d binhvuongos -f /path/000022_inbox_external_ref.down.sql'
```

## Risks
- Migration fail trên prod (hiếm — schema đơn giản) → có down migration
- API_KEY leak trong log nếu middleware log header → review log output trước/sau deploy
- n8n self-host chưa setup xong → test qua cURL là đủ cho Phase 5 acceptance, n8n có thể làm sau

## Unresolved questions
- n8n deploy ở server Docker `103.97.125.186` hay VPS khác? → user quyết sau deploy Go app
