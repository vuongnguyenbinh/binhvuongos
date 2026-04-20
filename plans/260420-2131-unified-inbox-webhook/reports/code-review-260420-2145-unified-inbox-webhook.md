# Code Review — Unified Inbox Webhook

**Scope:** 10 files, ~350 LOC. Focus: JSON+multipart ingest, idempotency, Drive upload.

## Critical

**1. Idempotency race (api_handlers.go:119-131, migrations/000022).** Check-then-insert without a tx or `ON CONFLICT`. Two concurrent requests with same `(source, external_ref)` both miss `GetInboxByExternalRef`, then one insert succeeds and the other returns 500 `DB_ERROR` (unique violation) instead of 200 duplicate. Fix (not applied): `INSERT ... ON CONFLICT (source, external_ref) WHERE external_ref IS NOT NULL DO NOTHING RETURNING ...`, then fall back to SELECT when no row returned.

**2. API key timing-safe compare (middleware/api_key.go:11).** `key != apiKey` leaks length/prefix via timing. Use `crypto/subtle.ConstantTimeCompare`. Low practical risk (single-tenant) but trivial fix.

**3. Orphan Drive files on failure (api_handlers.go:65-71 → 146-149).** If `uploadAttachmentToDrive` succeeds but DB insert later fails, the Drive file is orphaned with no compensating delete. Acceptable for Layer 1 but flag for tracking.

## High

**4. Duplicate `attachment_urls` in multipart (api_handlers.go:74-81).** Reads both `c.Context().QueryArgs().PeekMulti(...)` and `mf.Value["attachment_urls"]`. `QueryArgs()` is URL-query only (not multipart form) — line 74 is dead code in normal multipart but will double-count if URL query string is used. Drop line 74; `mf.Value` is the source of truth.

**5. `detectItemType` ignores attachments (inbox_webhook_helpers.go:34-39).** Multipart upload with empty `item_type` and non-URL content gets classified as `note`, not `image`/`file`. Should inspect mime from first attachment.

**6. Content-Type prefix check is fragile (api_handlers.go:54).** `strings.HasPrefix(contentType, "multipart/form-data")` works but misses case variance. Use `strings.EqualFold` on the parsed media type or `c.Is("multipart/form-data")`.

## Medium

**7. `isHTTPURL` indexing (inbox_webhook_helpers.go:43).** Byte indexing on potentially-UTF-8 content works only because ASCII prefix; fine but `strings.HasPrefix` reads clearer.

**8. `url` field missing in multipart path validation.** `url` is accepted from form but never passed into `validateInboxInput` (length bound only enforced on attachment URLs).

**9. Query file drift (query/inbox_items.sql vs inbox_items.sql.go).** `UpdateInboxItem` in `.sql` has 4 params (`content,url,item_type,company_id`); generated Go code has 3 (no `company_id`). Hand-written code diverges from "source" SQL — either delete the .sql file or keep it aligned.

**10. `GetInboxByExternalRef` query file missing column safeguard.** No `ORDER BY` with `LIMIT 1`; if partial unique index is dropped in future, non-deterministic row returned. Cosmetic given the unique index.

**11. `log.Fatalf` on missing owner (handler.go:29).** Acceptable fail-fast, but a single bad env var crashes the process on restart; consider logging warning and disabling `/api/v1/inbox` route instead.

## Low

**12. Error message leaks internals (api_handlers.go:19, 148).** `APIError("DB_ERROR", err.Error())` returns raw PG error text to clients. Log server-side, return generic message externally.

**13. `maxContentBytes = 10KB` uses byte-length on Go string — fine, but name suggests chars. Rename or document.

**14. Helper file name `inbox_webhook_helpers.go` OK; consider `inbox_webhook.go` and move `APICreateInbox` out of `api_handlers.go` for cohesion (YAGNI: leave if not growing).

## Positive

- Closed whitelist for `item_type` (good)
- Size caps on all string inputs
- Partial unique index correctly scoped `WHERE external_ref IS NOT NULL`
- Eager owner-user resolution avoids per-request lookup
- JSONB default `'[]'` prevents null rows

## Recommended Actions

1. Replace check-then-insert with `ON CONFLICT DO NOTHING` pattern (critical #1)
2. Timing-safe API key compare (critical #2)
3. Remove dead `QueryArgs().PeekMulti` branch (high #4)
4. Sync `query/inbox_items.sql` with generated code or delete (medium #9)
5. Use generic error message for DB errors (low #12)

## Unresolved Questions

- Is `query/*.sql` the source of truth or dead documentation? (generated code is hand-written per CLAUDE context "no sqlc")
- Should Drive orphan cleanup be deferred to a separate reaper task or handled inline?
- Max body size for JSON path — is there a global Fiber `BodyLimit`? Not visible in `fiber.Config`.
