package protocol

import "testing"

func TestTusFileIDFromPath(t *testing.T) {
	fileID, err := tusFileIDFromPath(
		"/api/ui/files/uploads/33333333-3333-4333-8333-333333333333%2Bmultipart-id",
		"/api/ui/files/uploads",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fileID != "33333333-3333-4333-8333-333333333333" {
		t.Fatalf("unexpected file id: %s", fileID)
	}
}

func TestTusFileIDFromPathRejectsMissingID(t *testing.T) {
	if _, err := tusFileIDFromPath("/api/ui/files/uploads", "/api/ui/files/uploads"); err == nil {
		t.Fatal("expected missing upload id to fail")
	}
}
