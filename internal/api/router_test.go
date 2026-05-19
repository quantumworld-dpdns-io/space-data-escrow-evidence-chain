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

func withAdmin(req *http.Request) {
	req.Header.Set("X-API-Key", "test-key")
}

func createEvidence(t *testing.T, r *Router, body map[string]any, idem string) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/evidence", bytes.NewReader(b))
	withAdmin(req)
	req.Header.Set("Content-Type", "application/json")
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	r.Handler().ServeHTTP(w, req)
	return w
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

func TestRBACViewerForbiddenOnAdminEndpoint(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/admin/key-rotation", nil)
	req.Header.Set("Authorization", "Bearer jwt-viewer")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", w.Code)
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

func TestBulkVerify(t *testing.T) {
	r := newTestRouter()
	created := createEvidence(t, r, map[string]any{"external_id": "EXT-B1", "source": "sat-b", "type": "imagery", "payload": map[string]string{"d": "x"}}, "")
	var rec map[string]any
	_ = json.Unmarshal(created.Body.Bytes(), &rec)
	id := rec["id"].(string)

	body, _ := json.Marshal(map[string]any{"ids": []string{id, "missing"}})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/verify/bulk", bytes.NewReader(body))
	withAdmin(req)
	req.Header.Set("Content-Type", "application/json")
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestAttestChainProofAndAdmin(t *testing.T) {
	r := newTestRouter()
	created := createEvidence(t, r, map[string]any{"external_id": "EXT-A1", "source": "sat-a", "type": "imagery", "payload": map[string]string{"d": "x"}}, "")
	var rec map[string]any
	_ = json.Unmarshal(created.Body.Bytes(), &rec)
	id := rec["id"].(string)

	attBody, _ := json.Marshal(map[string]any{"evidence_id": id, "signer": "sig-1", "signature": "abc", "algorithm": "ed25519"})
	wAtt := httptest.NewRecorder()
	reqAtt, _ := http.NewRequest(http.MethodPost, "/v1/attest", bytes.NewReader(attBody))
	withAdmin(reqAtt)
	reqAtt.Header.Set("Content-Type", "application/json")
	r.Handler().ServeHTTP(wAtt, reqAtt)
	if wAtt.Code != http.StatusAccepted {
		t.Fatalf("expected 202 got %d body=%s", wAtt.Code, wAtt.Body.String())
	}

	wChain := httptest.NewRecorder()
	reqChain, _ := http.NewRequest(http.MethodGet, "/v1/chain/"+id, nil)
	withAdmin(reqChain)
	r.Handler().ServeHTTP(wChain, reqChain)
	if wChain.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", wChain.Code)
	}

	wProof := httptest.NewRecorder()
	reqProof, _ := http.NewRequest(http.MethodGet, "/v1/proof/"+id, nil)
	withAdmin(reqProof)
	r.Handler().ServeHTTP(wProof, reqProof)
	if wProof.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", wProof.Code)
	}

	wAdmin := httptest.NewRecorder()
	reqAdmin, _ := http.NewRequest(http.MethodGet, "/v1/admin/key-rotation", nil)
	withAdmin(reqAdmin)
	r.Handler().ServeHTTP(wAdmin, reqAdmin)
	if wAdmin.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", wAdmin.Code)
	}
}

func TestSemanticSearchAndEnrichment(t *testing.T) {
	r := newTestRouter()
	_ = createEvidence(t, r, map[string]any{"external_id": "EXT-S1", "source": "sat-s", "type": "imagery", "payload": map[string]string{"desc": "storm over pacific"}}, "")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/search?mode=semantic&q=storm&page_size=5", nil)
	withAdmin(req)
	r.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
