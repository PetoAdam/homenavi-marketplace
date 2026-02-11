package handlers

import "testing"

func TestValidatePublishRequest(t *testing.T) {
	req := testPublishRequest(t)
	if err := validatePublishRequest(req); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}

	req.ID = ""
	if err := validatePublishRequest(req); err == nil {
		t.Fatalf("expected validation error for empty id")
	}
}

func TestValidatePublishRequestRejectsDevCompose(t *testing.T) {
	req := testPublishRequest(t)
	req.ComposeFile = newComposeServer(t, "services:\n  spotify:\n    image: ghcr.io/petoadam/homenavi-spotify:latest\n    volumes:\n      - ${HOMENAVI_ROOT}/integrations/secrets/spotify.secrets.json:/app/config/integration.secrets.json\n")
	if err := validatePublishRequest(req); err == nil {
		t.Fatalf("expected validation error for dev compose")
	}
}

func TestValidatePublishRequestRequiresImagePerService(t *testing.T) {
	req := testPublishRequest(t)
	req.ComposeFile = newComposeServer(t, "services:\n  spotify:\n    volumes:\n      - ${INTEGRATIONS_ROOT}/integrations/secrets/spotify.secrets.json:/app/config/integration.secrets.json\n")
	if err := validatePublishRequest(req); err == nil {
		t.Fatalf("expected validation error for missing image")
	}
}
