# CLAUDE.md — Bình Vương OS

## Project

Internal OS cho Solo Founder quản lý <10 công ty, ~20 nhân sự/CTV.

- **Stack:** Go (Fiber v2) + templ + Tailwind CSS v3 + HTMX + SortableJS
- **Status:** Layer 1 (Static UI demo, no DB)
- **Live:** https://os.binhvuong.vn
- **GitHub:** https://github.com/vuongnguyenbinh/binhvuongos

## Deploy Process

### Local → GitHub → Server

```bash
# 1. Commit
git add -A && git commit -m "feat/fix: description"

# 2. Push
git push

# 3. Deploy to server (SSH)
sshpass -p 'pi8phohVieyai9i' ssh root@103.97.125.186 \
  "cd /opt/binhvuongos && git pull && \
   docker rmi binhvuongos-app 2>/dev/null; \
   docker build --no-cache -t binhvuongos-app . && \
   docker compose up -d"
```

### Cache Busting

CSS/JS links include `?v=YYYYMMDD`. Update version in `web/templates/layout.templ` when making CSS/JS changes:
```
/static/css/output.css?v=20260416
/static/js/app.js?v=20260416
```

### Server Info

| Server | IP | User | Password | Purpose |
|--------|-----|------|----------|---------|
| Docker/App | 103.97.125.186 | root | pi8phohVieyai9i | Go app + Caddy + Cloudflare Tunnel |
| PostgreSQL | 103.97.125.131 | root | thochah5auch6Is | Database (Layer 2) |

### Cloudflare Tunnel

- Tunnel name: `binhvuong-os`
- Tunnel ID: `a4048ebd-132d-4135-a1d1-76998f8ee4a7`
- Domain: `os.binhvuong.vn` → `localhost:3000`
- Config: `/root/.cloudflared/config.yml` on server
- Service: `systemctl status cloudflared`

### Known Issues

- Docker build uses `--no-cache` to avoid stale Tailwind CSS output
- Cloudflare may cache CSS — always bump `?v=` version param
- `docker compose build` alone may use cached layers — use `docker rmi` + `docker build --no-cache` for guaranteed fresh build

## Architecture

```
cmd/server/main.go          → Fiber routes
internal/handler/*.go       → Route handlers (render templ)
web/templates/layout.templ  → Base layout (header, tabs, footer)
web/templates/pages/*.templ → Page templates
web/templates/components/   → Shared components (date_filter, pagination, modal)
web/static/css/input.css    → Tailwind input + dark mode CSS vars
web/static/js/              → app.js, dark-mode.js
tailwind.config.js          → Colors use CSS variables for auto dark/light
Dockerfile                  → Multi-stage: Go build + Tailwind build + Alpine runtime
```

## Design System

- **Fonts:** Phudu (display), Inter (body), JetBrains Mono (code/data)
- **Colors:** CSS-variable-based, auto-switch light↔dark via `.dark` class
- **Dark mode:** Tailwind `darkMode: 'class'`, localStorage persist, system preference detect

### Color Variables (light → dark)

| Variable | Light | Dark |
|----------|-------|------|
| `--ink` | #1A1918 | #F2F0ED |
| `--bg` | #FAF8F3 | #0D0F14 |
| `--surface` | #FFFFFF | #161920 |
| `--forest` | #1F3D2E | #7FD89E |
| `--muted` | #6B665E | #B0AAA0 |
| `--ember` | #D94F30 | #F07050 |
| `--sage` | #4A7C59 | #7FD89E |
| `--rust` | #A64545 | #F08080 |

## Routes

### List Pages (9 tabs)
- `/` — Dashboard
- `/inbox` — Hộp thư đến
- `/work-logs` — Nhật ký công việc
- `/tasks` — Công việc (Kanban)
- `/content` — Nội dung Pipeline
- `/companies` — Công ty
- `/campaigns` — Chiến dịch
- `/knowledge` — Kho kiến thức
- `/bookmarks` — Bookmark

### Detail Pages
- `/{module}/:id` — Content, Companies, Tasks, Work-logs, Campaigns, Knowledge
