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

func TestHealth(t *testing.T) {
	r := NewRouter(service.New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo()), "test-key")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestAuthRequired(t *testing.T) {
	r := NewRouter(service.New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo()), "test-key")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/search", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestCreateEvidence(t *testing.T) {
	r := NewRouter(service.New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo()), "test-key")
	payload := map[string]any{
		"external_id": "EXT-2",
		"source":      "sat-b",
		"type":        "telemetry",
		"payload":     map[string]string{"k": "v"},
	}
	b, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/evidence", bytes.NewReader(b))
	req.Header.Set("X-API-Key", "test-key")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
	}
}
