package middleware

import (
	"net/http"
	"strings"
)

type CORS struct {
	AllowedOrigins []string
}

func (c CORS) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if c.isAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Marketplace-Token")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c CORS) isAllowed(origin string) bool {
	if origin == "" || len(c.AllowedOrigins) == 0 {
		return false
	}
	for _, allowed := range c.AllowedOrigins {
		if strings.EqualFold(strings.TrimSpace(allowed), origin) {
			return true
		}
	}
	return false
}
