package auth

import (
	"net/http"
	"testing"
)

func TestAuthenticateAPIKey(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "k")
	p, err := Authenticate(req, "k")
	if err != nil || p.Role != "admin" {
		t.Fatalf("expected admin via api key, got p=%+v err=%v", p, err)
	}
}

func TestAuthenticateJWTViewer(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer jwt-viewer")
	p, err := Authenticate(req, "")
	if err != nil || p.Role != "viewer" {
		t.Fatalf("expected viewer via jwt, got p=%+v err=%v", p, err)
	}
}

func TestRBAC(t *testing.T) {
	if IsAllowed("viewer", "/v1/admin/key-rotation", http.MethodGet) {
		t.Fatal("viewer must not access admin path")
	}
	if !IsAllowed("admin", "/v1/admin/key-rotation", http.MethodGet) {
		t.Fatal("admin should access admin path")
	}
}
