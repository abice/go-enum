package middleware

import (
	"net/http"
)

type Middleware interface {
	Apply(http.Handler) http.Handler
}

// MiddlewareFunc implements Middleware interface
type MiddlewareFunc func(http.Handler) http.Handler

func (m MiddlewareFunc) Apply(h http.Handler) http.Handler {
	return m(h)
}

// RateLimiter
type rateLimiter struct {
	records map[string]int
}
