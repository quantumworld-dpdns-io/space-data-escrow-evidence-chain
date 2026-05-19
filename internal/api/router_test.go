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

func createEvidence(t *testing.T, r *Router, body map[string]any, idem string) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/evidence", bytes.NewReader(b))
	req.Header.Set("X-API-Key", "test-key")
	req.Header.Set("Content-Type", "application/json")
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	r.Handler().ServeHTTP(w, req)
	return w
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
	w := createEvidence(t, r, payload, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateEvidenceIdempotent(t *testing.T) {
	r := newTestRouter()
	payload := map[string]any{"external_id": "EXT-IDEMP", "source": "sat-id", "type": "imagery", "payload": map[string]string{"k": "v"}}
	w1 := createEvidence(t, r, payload, "idem-1")
	if w1.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", w1.Code)
	}
	w2 := createEvidence(t, r, payload, "idem-1")
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 for idempotent replay got %d", w2.Code)
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

func TestSearchPaginationShape(t *testing.T) {
	r := newTestRouter()
	for i := 0; i < 3; i++ {
		w := createEvidence(t, r, map[string]any{"external_id": "EXT-P", "source": "sat-p", "type": "imagery", "payload": map[string]string{"i": "x"}}, "")
		if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest {
			t.Fatalf("unexpected create code: %d", w.Code)
		}
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/search?page=1&page_size=2&source=sat-p", nil)
	req.Header.Set("X-API-Key", "test-key")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	var body ListEvidenceResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.Page != 1 || body.PageSize != 2 {
		t.Fatalf("unexpected pagination response: %+v", body)
	}
}
