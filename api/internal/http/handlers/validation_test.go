package handlers

import "testing"

func TestValidatePublishRequest(t *testing.T) {
	req := testPublishRequest()
	if err := validatePublishRequest(req); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}

	req.ID = ""
	if err := validatePublishRequest(req); err == nil {
		t.Fatalf("expected validation error for empty id")
	}
}
