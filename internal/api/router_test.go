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

func TestSearchPaginationShape(t *testing.T) {
	r := newTestRouter()
	for i := 0; i < 3; i++ {
		_ = createEvidence(t, r, map[string]any{"external_id": "EXT-P" + string(rune('A'+i)), "source": "sat-p", "type": "imagery", "payload": map[string]string{"i": "x"}}, "")
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

func TestSemanticSearch(t *testing.T) {
	r := newTestRouter()
	_ = createEvidence(t, r, map[string]any{"external_id": "EXT-S1", "source": "sat-s", "type": "imagery", "payload": map[string]string{"desc": "storm over pacific"}}, "")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/search?mode=semantic&q=storm&page_size=5", nil)
	req.Header.Set("X-API-Key", "test-key")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	var body ListEvidenceResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.Total < 1 {
		t.Fatalf("expected semantic results, got %+v", body)
	}
}

func TestEnrichmentTriggerAndStatus(t *testing.T) {
	r := newTestRouter()
	created := createEvidence(t, r, map[string]any{"external_id": "EXT-E1", "source": "sat-e", "type": "imagery", "payload": map[string]string{"desc": "launch event"}}, "")
	if created.Code != http.StatusCreated {
		t.Fatalf("expected create 201 got %d", created.Code)
	}
	var rec map[string]any
	_ = json.Unmarshal(created.Body.Bytes(), &rec)
	id := rec["id"].(string)

	enrichBody, _ := json.Marshal(TriggerEnrichmentRequest{EvidenceID: id})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/enrich", bytes.NewReader(enrichBody))
	req.Header.Set("X-API-Key", "test-key")
	req.Header.Set("Content-Type", "application/json")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202 got %d body=%s", w.Code, w.Body.String())
	}
	var job EnrichmentJobResponse
	_ = json.Unmarshal(w.Body.Bytes(), &job)
	if job.ID == "" || job.Status == "" {
		t.Fatalf("invalid job response: %+v", job)
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/v1/enrich/"+job.ID, nil)
	req2.Header.Set("X-API-Key", "test-key")
	r.Handler().ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w2.Code)
	}
}
