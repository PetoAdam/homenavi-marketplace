package server

import (
	"net/http"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/config"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/http/handlers"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.Config, pool *pgxpool.Pool) http.Handler {
	verifier := handlers.NewGitHubOIDCVerifier(cfg)
	return NewWithVerifier(cfg, pool, verifier)
}

func NewWithVerifier(cfg config.Config, pool *pgxpool.Pool, verifier handlers.OIDCVerifier) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.CORS{AllowedOrigins: cfg.AllowedOrigin}.Handler)

	h := handlers.IntegrationsHandler{Pool: pool, OIDCVerifier: verifier, OIDCTagPrefix: cfg.OIDCTagPrefix}

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api/integrations", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/publish-oidc", h.PublishOIDC)
		r.Get("/{id}", h.Get)
		r.Get("/{id}/versions", h.Versions)
	})

	return r
}
