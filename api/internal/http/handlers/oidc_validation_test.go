package handlers

import "testing"

func TestValidateOIDCRequest(t *testing.T) {
	claims := OIDCClaims{
		Repository: "PetoAdam/homenavi-spotify",
		Ref:        "refs/tags/v0.1.0",
		RefType:    "tag",
		SHA:        "abc123",
	}

	req := testPublishRequest(t)
	req.Version = "v0.1.0"
	req.ReleaseTag = "v0.1.0"
	req.RepoURL = "https://github.com/PetoAdam/homenavi-spotify"
	req.ManifestURL = "https://raw.githubusercontent.com/PetoAdam/homenavi-spotify/v0.1.0/manifest/homenavi-integration.json"

	if err := validateOIDCRequest(req, claims, "v0.1.0"); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}
}

func TestValidateOIDCRequestRejectsRepoMismatch(t *testing.T) {
	claims := OIDCClaims{
		Repository: "PetoAdam/homenavi-spotify",
		Ref:        "refs/tags/v0.1.0",
		RefType:    "tag",
		SHA:        "abc123",
	}

	req := testPublishRequest(t)
	req.Version = "v0.1.0"
	req.ReleaseTag = "v0.1.0"
	req.RepoURL = "https://github.com/PetoAdam/other-repo"
	req.ManifestURL = "https://raw.githubusercontent.com/PetoAdam/homenavi-spotify/v0.1.0/manifest/homenavi-integration.json"

	if err := validateOIDCRequest(req, claims, "v0.1.0"); err == nil {
		t.Fatalf("expected repo mismatch error")
	}
}

func TestTagFromClaims(t *testing.T) {
	claims := OIDCClaims{
		Ref:     "refs/tags/v1.2.3",
		RefType: "tag",
	}

	tag, err := tagFromClaims(claims, "v")
	if err != nil {
		t.Fatalf("expected tag, got %v", err)
	}
	if tag != "v1.2.3" {
		t.Fatalf("expected v1.2.3, got %s", tag)
	}
}
