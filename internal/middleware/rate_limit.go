package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type loginAttempt struct {
	count    int
	lockedAt time.Time
}

var (
	attempts = make(map[string]*loginAttempt)
	mu       sync.Mutex
)

const (
	maxAttempts = 5
	lockDuration = 15 * time.Minute
)

// LoginRateLimit blocks login after 5 failed attempts per IP for 15 minutes
func LoginRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		mu.Lock()
		a, exists := attempts[ip]
		if exists && !a.lockedAt.IsZero() && time.Since(a.lockedAt) < lockDuration {
			mu.Unlock()
			return c.Status(429).SendString("Quá nhiều lần thử. Vui lòng đợi 15 phút.")
		}
		if exists && !a.lockedAt.IsZero() && time.Since(a.lockedAt) >= lockDuration {
			// Reset after lock duration
			a.count = 0
			a.lockedAt = time.Time{}
		}
		mu.Unlock()
		return c.Next()
	}
}

// RecordFailedLogin increments failed attempt count for IP
func RecordFailedLogin(ip string) {
	mu.Lock()
	defer mu.Unlock()
	a, exists := attempts[ip]
	if !exists {
		a = &loginAttempt{}
		attempts[ip] = a
	}
	a.count++
	if a.count >= maxAttempts {
		a.lockedAt = time.Now()
	}
}

// ClearLoginAttempts resets attempts on successful login
func ClearLoginAttempts(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(attempts, ip)
}
