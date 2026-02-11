package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/models"
)

func testComposeYAML() string {
	return "services:\n  spotify:\n    image: ghcr.io/petoadam/homenavi-spotify:latest\n    volumes:\n      - ${INTEGRATIONS_ROOT}/integrations/secrets/spotify.secrets.json:/app/config/integration.secrets.json\n"
}

func newComposeServer(t *testing.T, composeYAML string) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/compose/docker-compose.integration.yml" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/yaml")
		_, _ = w.Write([]byte(composeYAML))
	}))
	t.Cleanup(server.Close)
	return server.URL + "/compose/docker-compose.integration.yml"
}

func testPublishRequest(t *testing.T) models.PublishRequest {
	composeURL := newComposeServer(t, testComposeYAML())
	return models.PublishRequest{
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
		ComposeFile: composeURL,
	}
}
