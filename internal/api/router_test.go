package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

func newTestRouter() *Router {
	return NewRouter(
		service.New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo()),
		"test-key",
		map[string]string{"version": "test", "commit": "abc123", "build_date": "2026-05-19"},
	)
}

func TestHealth(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestVersion(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/version", nil)
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["commit"] != "abc123" {
		t.Fatalf("expected commit abc123 got %q", body["commit"])
	}
}

func TestAuthRequired(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/search", nil)
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestCreateEvidence(t *testing.T) {
	r := newTestRouter()
	payload := map[string]any{"external_id": "EXT-2", "source": "sat-b", "type": "telemetry", "payload": map[string]string{"k": "v"}}
	b, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/evidence", bytes.NewReader(b))
	req.Header.Set("X-API-Key", "test-key")
	req.Header.Set("Content-Type", "application/json")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCustodyAndVerifyNotFound(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/verify/not-exists", nil)
	req.Header.Set("X-API-Key", "test-key")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuditEndpoint(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/audit", nil)
	req.Header.Set("X-API-Key", "test-key")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	var body AuditResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.Entries == nil {
		t.Fatal("expected entries array")
	}
}
