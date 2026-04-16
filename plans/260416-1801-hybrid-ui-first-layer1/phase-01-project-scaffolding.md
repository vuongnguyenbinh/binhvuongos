# Phase 1: Project Scaffolding & Dev Tooling

## Overview
- **Priority:** P0 — blocks all other phases
- **Status:** Pending
- **Effort:** 3h
- Setup Go project, Fiber server, templ, Tailwind CSS build, Air hot reload, Docker dev env

## Context Links
- PRD: `Root/binhvuong-os-prd.md` Section 1.3 (stack)
- Demo: `Root/binhvuong-os-demo.html` (reference for Tailwind config)

## Key Insights
- templ generates Go code from `.templ` files → needs `templ generate` step
- Tailwind CSS needs Node.js CLI for build (not CDN in production)
- Air watches `.go` and `.templ` files for hot reload
- Fiber v2 serves static files + templ-rendered HTML

## Requirements

### Functional
- `make dev` → starts Air + Tailwind watch + templ watch
- `make build` → compiles Go binary + generates CSS
- Browser hits `localhost:3000` → sees "Bình Vương OS" placeholder page
- Static files served from `/static/`

### Non-functional
- Hot reload < 2s for templ/CSS changes
- Docker Compose for consistent dev environment

## Files to Create

| File | Purpose |
|------|---------|
| `go.mod` | Go module: `binhvuongos` |
| `go.sum` | Dependencies |
| `cmd/server/main.go` | Fiber entry point, routes, static file serving |
| `web/tailwind.config.js` | Tailwind config matching demo color palette (with updated colors) |
| `web/static/css/input.css` | Tailwind directives + custom CSS from demo |
| `web/static/js/app.js` | Tab switching + checkbox logic (from demo `<script>`) |
| `Makefile` | dev, build, templ-generate, tailwind-build targets |
| `.air.toml` | Air hot reload config |
| `docker-compose.yml` | Go app service (Postgres placeholder for Layer 2) |
| `Dockerfile` | Multi-stage: build + runtime |

## Implementation Steps

### 1. Init Go module
```bash
cd binhvuongos
go mod init binhvuongos
go get github.com/gofiber/fiber/v2
go get github.com/a-h/templ
```

### 2. Create Fiber server (`cmd/server/main.go`)
```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "binhvuongos/internal/handler"
)

func main() {
    app := fiber.New(fiber.Config{
        ViewsLayout: "",
    })
    app.Use(logger.New())
    app.Static("/static", "./web/static")

    // Pages — all return templ-rendered HTML
    app.Get("/", handler.Dashboard)
    app.Get("/work-logs", handler.WorkLogs)
    app.Get("/inbox", handler.Inbox)
    app.Get("/tasks", handler.Tasks)
    app.Get("/content", handler.Content)
    app.Get("/companies", handler.Companies)
    app.Get("/campaigns", handler.Campaigns)
    app.Get("/knowledge", handler.Knowledge)

    log.Fatal(app.Listen(":3000"))
}
```

### 3. Create placeholder handler (`internal/handler/dashboard.go`)
```go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "binhvuongos/web/templates/pages"
    "github.com/a-h/templ"
)

func Dashboard(c *fiber.Ctx) error {
    component := pages.DashboardPage()
    c.Set("Content-Type", "text/html")
    return component.Render(c.Context(), c.Response().BodyWriter())
}
```

### 4. Tailwind config (`web/tailwind.config.js`)
Extract from demo HTML, replace amber2/ochre with ember/flame:
```js
module.exports = {
  content: ["./web/templates/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Phudu', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      colors: {
        ivory: '#FAF8F3',
        surface: '#FFFFFF',
        ink: '#1A1918',
        forest: { DEFAULT: '#1F3D2E', 50: '#E8EDE9', 100: '#C4D1C8', 600: '#2D5C44', 900: '#122419' },
        ember: '#D94F30',      // was amber2 #B8741F
        flame: '#E8623A',      // was ochre #C89B3C
        muted: '#6B665E',
        hairline: '#E8E4DB',
        cream: '#F2EEE4',
        sage: '#4A7C59',
        rust: '#A64545',
      }
    }
  }
}
```

### 5. Input CSS (`web/static/css/input.css`)
Port all custom CSS from demo `<style>` block:
- `.display`, `.mono`, `.tnum`, `.eyebrow`, `.ornament`
- `.pill`, `.progress-track`, `.progress-fill`
- `.kanban-card`, `.stat-hero`, `.check`, `.spark-bar`
- `.tab-active`, `.row-hover`, `.paper`
- Scrollbar styles, animations (`fadeIn`)
- Replace all `--amber` refs with `--ember`

### 6. Makefile
```makefile
.PHONY: dev build templ-generate tailwind-build

dev:
	@echo "Starting dev server..."
	@make -j3 templ-watch tailwind-watch air

templ-watch:
	templ generate --watch

templ-generate:
	templ generate

tailwind-watch:
	npx tailwindcss -i web/static/css/input.css -o web/static/css/output.css --watch

tailwind-build:
	npx tailwindcss -i web/static/css/input.css -o web/static/css/output.css --minify

air:
	air

build: templ-generate tailwind-build
	go build -o bin/server cmd/server/main.go
```

### 7. Air config (`.air.toml`)
```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/server"
  bin = "./tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "web/static", "Root", "plans"]

[misc]
  clean_on_exit = true
```

### 8. Docker Compose (`docker-compose.yml`)
```yaml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    volumes:
      - .:/app
    working_dir: /app
    command: make dev
```

### 9. Dockerfile (multi-stage)
```dockerfile
FROM golang:1.22-alpine AS builder
RUN go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN templ generate && go build -o bin/server cmd/server/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/web/static ./web/static
EXPOSE 3000
CMD ["./server"]
```

## Todo List
- [ ] Init Go module + install deps (fiber, templ)
- [ ] Create `cmd/server/main.go` with routes
- [ ] Create placeholder handler for dashboard
- [ ] Setup Tailwind config with updated colors
- [ ] Port custom CSS from demo to `input.css`
- [ ] Port JS (tab switching, checkbox) to `app.js`
- [ ] Create Makefile with dev/build targets
- [ ] Create Air config
- [ ] Create Docker Compose + Dockerfile
- [ ] Verify `make dev` → browser shows placeholder page

## Success Criteria
- `make dev` starts without errors
- `localhost:3000` returns HTML with correct fonts + colors
- Tailwind classes compile correctly
- templ generates Go code from `.templ` files
- Hot reload works (edit templ → browser refreshes)

## Risk Assessment
- **templ + Fiber integration**: templ outputs `io.Writer`, Fiber uses `*fiber.Ctx` — need adapter pattern. Mitigation: use `component.Render(ctx, writer)` pattern.
- **Tailwind content path**: Must scan `.templ` files, not `.html`. Verify content config.

## Next Steps
→ Phase 2: Design system components (shared templ components)
