package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/config"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/http/handlers"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/http/server"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/models"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/testutil"
)

type stubOIDCVerifier struct {
	claims      handlers.OIDCClaims
	verifyErr   error
	workflowErr error
}

func (s stubOIDCVerifier) Verify(_ context.Context, _ string) (handlers.OIDCClaims, error) {
	return s.claims, s.verifyErr
}

func (s stubOIDCVerifier) VerifyWorkflow(_ context.Context, _ handlers.OIDCClaims) error {
	return s.workflowErr
}

func TestPublishOIDCRequiresToken(t *testing.T) {
	pool, cleanup := testutil.StartPostgres(t)
	defer cleanup()

	verifier := stubOIDCVerifier{}
	h := server.NewWithVerifier(config.Config{}, pool, verifier)

	reqBody := models.PublishRequest{
		ID:          "spotify",
		Name:        "Spotify",
		Version:     "v0.1.0",
		Description: "Play music",
		ManifestURL: "https://example.com/manifest.json",
		Manifest:    map[string]any{"id": "spotify"},
		Image:       "ghcr.io/petoadam/homenavi-spotify:latest",
		Images:      []string{},
		Assets:      map[string]string{},
		ListenPath:  "/integrations/spotify",
	}
	payload, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/publish-oidc", bytes.NewReader(payload))
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
}

func TestPublishOIDCAndList(t *testing.T) {
	pool, cleanup := testutil.StartPostgres(t)
	defer cleanup()

	verifier := stubOIDCVerifier{claims: handlers.OIDCClaims{
		Repository: "PetoAdam/homenavi-spotify",
		Ref:        "refs/tags/v0.1.0",
		RefType:    "tag",
		SHA:        "abc123",
	}}
	h := server.NewWithVerifier(config.Config{OIDCTagPrefix: "v"}, pool, verifier)

	reqBody := models.PublishRequest{
		ID:          "spotify",
		Name:        "Spotify",
		Version:     "v0.1.0",
		Description: "Play music",
		ManifestURL: "https://raw.githubusercontent.com/PetoAdam/homenavi-spotify/v0.1.0/manifest/homenavi-integration.json",
		Manifest:    map[string]any{"id": "spotify"},
		Image:       "ghcr.io/petoadam/homenavi-spotify:latest",
		Images:      []string{},
		Assets:      map[string]string{},
		ListenPath:  "/integrations/spotify",
		RepoURL:     "https://github.com/PetoAdam/homenavi-spotify",
		ReleaseTag:  "v0.1.0",
	}
	payload, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/publish-oidc", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer test-token")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/integrations", nil)
	listRes := httptest.NewRecorder()
	h.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d", listRes.Code)
	}
}
