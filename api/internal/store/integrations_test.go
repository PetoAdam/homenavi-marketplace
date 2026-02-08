package store

import (
	"context"
	"testing"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/models"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/testutil"
)

func TestPublishIntegrationAndLatest(t *testing.T) {
	pool, cleanup := testutil.StartPostgres(t)
	defer cleanup()

	ctx := context.Background()
	baseReq := models.PublishRequest{
		ID:          "spotify",
		Name:        "Spotify",
		Version:     "v0.1.0",
		Description: "Play music",
		ManifestURL: "https://example.com/manifest.json",
		Manifest:    map[string]any{"id": "spotify"},
		Image:       "ghcr.io/petoadam/homenavi-spotify:latest",
		Images:      []string{"https://example.com/hero.png"},
		Assets:      map[string]string{"icon": "https://example.com/icon.svg"},
		ListenPath:  "/integrations/spotify",
		RepoURL:     "https://github.com/PetoAdam/homenavi-spotify",
		ReleaseTag:  "v0.1.0",
		Publisher:   "Homenavi",
	}

	item, err := PublishIntegration(ctx, pool, baseReq, true)
	if err != nil {
		t.Fatalf("publish v0.1.0: %v", err)
	}
	if !item.Latest {
		t.Fatalf("expected latest true")
	}

	baseReq.Version = "v0.2.0"
	baseReq.ReleaseTag = "v0.2.0"
	item2, err := PublishIntegration(ctx, pool, baseReq, true)
	if err != nil {
		t.Fatalf("publish v0.2.0: %v", err)
	}
	if !item2.Latest {
		t.Fatalf("expected latest true after second publish")
	}

	latest, err := GetIntegration(ctx, pool, "spotify", "")
	if err != nil {
		t.Fatalf("get latest: %v", err)
	}
	if latest.Version != "v0.2.0" {
		t.Fatalf("expected latest version v0.2.0, got %s", latest.Version)
	}

	versions, err := ListVersions(ctx, pool, "spotify")
	if err != nil {
		t.Fatalf("list versions: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
}

func TestListenPathUnique(t *testing.T) {
	pool, cleanup := testutil.StartPostgres(t)
	defer cleanup()

	ctx := context.Background()
	req := models.PublishRequest{
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
		RepoURL:     "https://github.com/PetoAdam/homenavi-spotify",
		ReleaseTag:  "v0.1.0",
		Publisher:   "Homenavi",
	}

	if _, err := PublishIntegration(ctx, pool, req, true); err != nil {
		t.Fatalf("publish initial: %v", err)
	}

	req.ID = "alt-spotify"
	req.Version = "v0.1.0"
	if _, err := PublishIntegration(ctx, pool, req, true); err == nil {
		t.Fatalf("expected listen_path conflict")
	}
}
