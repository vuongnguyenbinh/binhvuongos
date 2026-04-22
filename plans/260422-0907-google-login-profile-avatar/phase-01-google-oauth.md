# Phase 1 — Google OAuth Login

**Effort:** 75m

## Files

### New
- `internal/oauth/google.go` — OAuth client helpers
- `internal/handler/google_auth.go` — GoogleLoginRedirect + GoogleCallback

### Modify
- `internal/config/config.go` — `GoogleRedirectURI`
- `cmd/server/main.go` — 2 routes
- `web/templates/pages/login.templ` — Google button

## oauth/google.go

```go
package oauth

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

type GoogleUserinfo struct {
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
    Name          string `json:"name"`
    Picture       string `json:"picture"`
}

type Config struct {
    ClientID     string
    ClientSecret string
    RedirectURI  string
}

// AuthorizeURL builds the Google consent page URL for the given state.
func (c *Config) AuthorizeURL(state string) string {
    q := url.Values{
        "client_id":     {c.ClientID},
        "redirect_uri":  {c.RedirectURI},
        "response_type": {"code"},
        "scope":         {"openid email profile"},
        "state":         {state},
        "access_type":   {"online"},
        "prompt":        {"select_account"},
    }
    return "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
}

// ExchangeCode swaps the auth code for an access token.
func (c *Config) ExchangeCode(ctx context.Context, code string) (string, error) {
    body := url.Values{
        "code":          {code},
        "client_id":     {c.ClientID},
        "client_secret": {c.ClientSecret},
        "redirect_uri":  {c.RedirectURI},
        "grant_type":    {"authorization_code"},
    }
    req, _ := http.NewRequestWithContext(ctx, "POST",
        "https://oauth2.googleapis.com/token",
        strings.NewReader(body.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := http.DefaultClient.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    var out struct {
        AccessToken string `json:"access_token"`
        ErrorDesc   string `json:"error_description"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if out.AccessToken == "" { return "", fmt.Errorf("oauth: %s", out.ErrorDesc) }
    return out.AccessToken, nil
}

// FetchUserinfo calls OpenID Connect userinfo endpoint.
func (c *Config) FetchUserinfo(ctx context.Context, accessToken string) (*GoogleUserinfo, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET",
        "https://openidconnect.googleapis.com/v1/userinfo", nil)
    req.Header.Set("Authorization", "Bearer "+accessToken)
    resp, err := http.DefaultClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    var u GoogleUserinfo
    if err := json.NewDecoder(resp.Body).Decode(&u); err != nil { return nil, err }
    return &u, nil
}
```

## handler/google_auth.go

```go
package handler

import (
    "crypto/rand"
    "encoding/hex"
    "strings"
    "time"

    "binhvuongos/internal/middleware"
    "binhvuongos/internal/oauth"

    "github.com/gofiber/fiber/v2"
)

func (h *Handler) newGoogleOAuth() *oauth.Config {
    return &oauth.Config{
        ClientID:     h.config.GoogleClientID,
        ClientSecret: h.config.GoogleClientSecret,
        RedirectURI:  h.config.GoogleRedirectURI,
    }
}

// GoogleLoginRedirect generates a CSRF state cookie and sends the user to Google.
func (h *Handler) GoogleLoginRedirect(c *fiber.Ctx) error {
    if h.config.GoogleClientID == "" || h.config.GoogleRedirectURI == "" {
        return c.Status(503).SendString("Google OAuth chưa được cấu hình")
    }
    buf := make([]byte, 16)
    _, _ = rand.Read(buf)
    state := hex.EncodeToString(buf)
    c.Cookie(&fiber.Cookie{
        Name:     "gstate",
        Value:    state,
        Path:     "/",
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Lax",
        Expires:  time.Now().Add(10 * time.Minute),
    })
    return c.Redirect(h.newGoogleOAuth().AuthorizeURL(state))
}

// GoogleCallback exchanges code, fetches email, looks up user, issues JWT.
func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
    state := c.Query("state")
    if state == "" || state != c.Cookies("gstate") {
        return c.Status(400).SendString("State không hợp lệ (có thể link đã hết hạn)")
    }
    c.ClearCookie("gstate")

    code := c.Query("code")
    if code == "" {
        return c.Status(400).SendString("Thiếu code từ Google")
    }
    oc := h.newGoogleOAuth()
    token, err := oc.ExchangeCode(c.Context(), code)
    if err != nil {
        return c.Status(502).SendString("Google exchange fail: " + err.Error())
    }
    info, err := oc.FetchUserinfo(c.Context(), token)
    if err != nil || info.Email == "" {
        return c.Status(502).SendString("Không lấy được email từ Google")
    }
    if !info.EmailVerified {
        return c.Status(403).SendString("Email chưa verify ở Google")
    }
    email := strings.ToLower(strings.TrimSpace(info.Email))

    user, err := h.queries.GetUserByEmail(c.Context(), email)
    if err != nil {
        return c.Status(403).SendString("Tài khoản " + email + " chưa được cấp quyền. Liên hệ admin.")
    }
    if user.Status != "active" {
        return c.Status(403).SendString("Tài khoản đã bị khoá")
    }
    // Populate avatar from Google picture if empty
    if (user.AvatarURL.String == "" || !user.AvatarURL.Valid) && info.Picture != "" {
        _ = h.queries.UpdateUserAvatar(c.Context(), user.ID, info.Picture)
    }
    _ = h.queries.UpdateLastLogin(c.Context(), user.ID)

    jwtToken, err := middleware.GenerateToken(user, h.config.JWTSecret, false)
    if err != nil {
        return c.Status(500).SendString("Không tạo được JWT")
    }
    c.Cookie(&fiber.Cookie{
        Name:     "token",
        Value:    jwtToken,
        Path:     "/",
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Lax",
        Expires:  time.Now().Add(24 * time.Hour),
    })
    return c.Redirect("/inbox")
}
```

## Login page button

Login.templ — append below form:
```html
<div class="mt-4 text-center">
    <a href="/auth/google" class="inline-flex items-center gap-2 px-4 py-2.5 border border-hairline rounded font-medium hover:bg-cream/30">
        <svg width="16" height="16" viewBox="0 0 48 48">... (Google G logo) ...</svg>
        Đăng nhập với Google
    </a>
</div>
```

## Routes

```go
// main.go — BEFORE AuthRequired group:
app.Get("/auth/google", h.GoogleLoginRedirect)
app.Get("/auth/google/callback", h.GoogleCallback)
```

## Queries new

Add `UpdateUserAvatar` in `users.sql.go`:
```go
func (q *Queries) UpdateUserAvatar(ctx context.Context, id pgtype.UUID, url string) error {
    _, err := q.pool.Exec(ctx, "UPDATE users SET avatar_url=$2 WHERE id=$1", id, url)
    return err
}
```

## Config

```go
// config.go
GoogleRedirectURI: getenvOr("GOOGLE_REDIRECT_URI", "https://os.binhvuong.vn/auth/google/callback"),
```

## Todo
- [ ] Create `internal/oauth/google.go`
- [ ] Create `internal/handler/google_auth.go`
- [ ] Add `UpdateUserAvatar` query
- [ ] Add `GoogleRedirectURI` to config
- [ ] Register 2 routes in main.go
- [ ] Add Google button to login.templ
- [ ] Verify Google Cloud Console redirect URI configured (manual step by user)

## Success criteria
- Click Google button → consent → callback → /inbox (authenticated)
- Unknown email → 403 with clear message
- Inactive user email → 403 blocked
- Avatar auto-populated from Google picture when null
