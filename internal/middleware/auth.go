package middleware

import (
	"encoding/hex"
	"fmt"
	"time"

	"binhvuongos/internal/config"
	"binhvuongos/internal/db/generated"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for the given user
func GenerateToken(user generated.User, secret string, rememberMe bool) (string, error) {
	expiry := 24 * time.Hour
	if rememberMe {
		expiry = 30 * 24 * time.Hour
	}

	claims := Claims{
		UserID: UUIDToString(user.ID),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// AuthRequired middleware checks for valid JWT cookie, redirects to /login if invalid
func AuthRequired(queries *generated.Queries, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")
		if tokenStr == "" {
			return c.Redirect("/login")
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.ClearCookie("token")
			return c.Redirect("/login")
		}

		uid := StringToUUID(claims.UserID)
		user, err := queries.GetUserByID(c.Context(), uid)
		if err != nil {
			c.ClearCookie("token")
			return c.Redirect("/login")
		}

		c.Locals("user", user)
		return c.Next()
	}
}

// UUIDToString converts pgtype.UUID to string
func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// StringToUUID converts a UUID string to pgtype.UUID
func StringToUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	if len(s) != 36 {
		return u
	}
	// Remove hyphens
	clean := s[0:8] + s[9:13] + s[14:18] + s[19:23] + s[24:36]
	b, err := hex.DecodeString(clean)
	if err != nil || len(b) != 16 {
		return u
	}
	copy(u.Bytes[:], b)
	u.Valid = true
	return u
}
