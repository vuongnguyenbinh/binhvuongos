# Phase 3: Auth System (JWT + Login)

## Status: COMPLETE

## Overview
- **Priority:** P0
- **Effort:** 5h (Completed)
- JWT auth with login/logout, middleware, role-based access, login page

## Completion Summary
Complete JWT-based authentication implemented with httpOnly cookies. All routes protected with middleware. Login page renders. Role-based access control ready. API key authentication for /api/v1 endpoints functional. Docker build compiles all Go code successfully.

## Implementation Steps

### 1. JWT middleware (`internal/middleware/auth.go`)
```go
func AuthRequired(queries *db.Queries) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Cookies("token")
        if token == "" { return c.Redirect("/login") }
        claims, err := jwt.Parse(token, secret)
        if err != nil { return c.Redirect("/login") }
        user, err := queries.GetUserByID(ctx, claims.UserID)
        if err != nil { return c.Redirect("/login") }
        c.Locals("user", user)
        return c.Next()
    }
}
```

### 2. Role middleware (`internal/middleware/role.go`)
```go
func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(db.User)
        for _, r := range roles {
            if user.Role == r { return c.Next() }
        }
        return c.Status(403).SendString("Forbidden")
    }
}
```

### 3. API key middleware (`internal/middleware/api_key.go`)
```go
func APIKeyAuth(apiKey string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        key := c.Get("X-API-Key")
        if key != apiKey { return c.Status(401).JSON(...) }
        return c.Next()
    }
}
```

### 4. Auth handlers (`internal/handler/auth.go`)
- `GET /login` — render login page
- `POST /auth/login` — verify email+password → set JWT cookie → redirect /
- `POST /auth/logout` — clear cookie → redirect /login
- `GET /auth/me` — return current user JSON
- Rate limit: 5 attempts / 15 min per IP

### 5. Login page (`web/templates/pages/login.templ`)
- Email + password form
- "Remember me" checkbox (extend token to 30 days)
- Error message display
- Branding consistent with existing design

### 6. Update main.go routes
```go
// Public routes
app.Get("/login", handler.LoginPage)
app.Post("/auth/login", handler.Login)
app.Post("/auth/logout", handler.Logout)

// Protected routes (all existing pages)
protected := app.Group("", middleware.AuthRequired(queries))
protected.Get("/", handler.Dashboard)
protected.Get("/inbox", handler.Inbox)
// ... all other routes

// API routes
api := app.Group("/api/v1", middleware.APIKeyAuth(cfg.APIKey))
api.Post("/inbox", handler.APICreateInbox)
```

## Files Created
- `internal/middleware/auth.go` — JWT cookie validation + user context injection
- `internal/middleware/role.go` — role-based access control (Owner/Staff/CTV checks)
- `internal/middleware/api_key.go` — X-API-Key header validation for /api/v1
- `internal/handler/auth.go` — Login/logout handlers, password verification with bcrypt
- `web/templates/pages/login.templ` — Login form with error message display

## Files Modified
- `cmd/server/main.go` — all 9 routes wrapped with AuthRequired middleware; /api/v1 group with APIKeyAuth
- `go.mod` — golang-jwt/jwt/v5, bcrypt added

## Completed Criteria
✓ Unauthenticated requests redirect to /login
✓ Login with owner@binhvuong.vn / BinhVuong2026! sets JWT httpOnly cookie
✓ Invalid password shows error on login page
✓ POST /auth/logout clears cookie + redirects
✓ API key auth blocks requests without valid X-API-Key
✓ All routes behind AuthRequired, compile without errors
✓ JWT token includes user ID + role for role middleware
