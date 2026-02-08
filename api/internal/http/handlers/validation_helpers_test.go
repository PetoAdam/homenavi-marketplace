package handlers

import "github.com/PetoAdam/homenavi-marketplace/api/internal/models"

func testPublishRequest() models.PublishRequest {
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
	}
}
