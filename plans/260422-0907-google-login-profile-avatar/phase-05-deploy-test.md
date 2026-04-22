# Phase 5 — Deploy + E2E Test

**Effort:** 30m

## Pre-deploy checklist

- [ ] Google Cloud Console → OAuth 2.0 Client ID → Authorized redirect URIs có `https://os.binhvuong.vn/auth/google/callback`
- [ ] `.env` trên server có `GOOGLE_REDIRECT_URI` (optional; defaults to production URL)
- [ ] `go build ./...` + `templ generate` pass locally via Docker

## Deploy

```bash
git add -A
git commit -m "feat: google login + self-profile + avatar upload + user detail"
git push

sshpass -p 'pi8phohVieyai9i' ssh root@103.97.125.186 \
  "cd /opt/binhvuongos && git pull && \
   docker rmi binhvuongos-app 2>/dev/null; \
   docker build --no-cache -t binhvuongos-app . && \
   docker compose up -d"
```

## E2E tests

### A. Google login (happy path — use real Google account)
1. Logout if logged in
2. Click "Đăng nhập với Google" on /login
3. Google consent → grant
4. Redirect back → /inbox landed, JWT cookie issued

### B. Google login — email not whitelisted
1. Login with Google account whose email NOT in users DB
2. Expect 403 with message

### C. Google login — CSRF state protection
1. Manually hit `/auth/google/callback?code=fake&state=wrong` → 400

### D. Self-profile edit
1. Login, go to /profile
2. Change full_name + phone → save
3. Verify DB row updated; profile page reflects

### E. Avatar upload
1. Pick image <10MB → upload
2. Profile page shows image
3. Header shows image in avatar slot
4. Try >10MB image → 400 error

### F. Avatar auto-populate on first Google login
1. User with null avatar_url → login Google with profile picture
2. DB users.avatar_url should be set after callback

### G. User detail admin view
1. Owner visit /users/:staff_id → profile + tasks + logs render
2. Manager visit /users/:owner_id → 200 (view) but no edit button
3. Staff visit /users/:any_id → 403

### H. Header dropdown
1. Click avatar → menu appears
2. Click outside → menu closes
3. Profile link → /profile
4. Logout button → /login

### I. Regression
- Email/password login still works
- Webhook /api/v1/inbox still 201
- Triage, user CRUD unaffected

## Cleanup
No test data to remove (Google login uses real account; no synthetic rows).

## Todo
- [ ] User confirms Google Cloud Console redirect URI configured
- [ ] Deploy
- [ ] Verify startup log + build image
- [ ] Run E2E A–I
- [ ] Update backlog.md — close related items

## Success criteria
- All E2E groups pass
- Zero regression on existing auth paths
- Backlog reconciled

## Risks
- Redirect URI mismatch → Google error screen; user must add URI in Cloud Console first
- OAuth app in "Testing" mode on Google limits 100 users — bump to "Production" if needed (or add users as test users)
