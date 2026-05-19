package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

type Router struct {
	svc    *service.Service
	apiKey string
	meta   map[string]string
	mux    *http.ServeMux
}

func NewRouter(svc *service.Service, apiKey string, meta map[string]string) *Router {
	r := &Router{svc: svc, apiKey: apiKey, meta: meta, mux: http.NewServeMux()}
	r.routes()
	return r
}

func (r *Router) Handler() http.Handler { return r.mux }

func (r *Router) routes() {
	r.mux.HandleFunc("/healthz", r.health)
	r.mux.HandleFunc("/readyz", r.ready)
	r.mux.HandleFunc("/version", r.version)
	r.mux.HandleFunc("/v1/evidence", r.withAuth(r.evidenceCreate))
	r.mux.HandleFunc("/v1/evidence/", r.withAuth(r.evidenceGet))
	r.mux.HandleFunc("/v1/custody", r.withAuth(r.custodyAppend))
	r.mux.HandleFunc("/v1/verify/", r.withAuth(r.verify))
	r.mux.HandleFunc("/v1/search", r.withAuth(r.search))
	r.mux.HandleFunc("/v1/audit", r.withAuth(r.audit))
}

func (r *Router) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if r.apiKey != "" && req.Header.Get("X-API-Key") != r.apiKey {
			writeJSON(w, http.StatusUnauthorized, APIError{Error: "unauthorized", Code: "AUTH_001"})
			return
		}
		next(w, req)
	}
}

func (r *Router) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}

func (r *Router) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (r *Router) version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"version":    r.meta["version"],
		"commit":     r.meta["commit"],
		"build_date": r.meta["build_date"],
	})
}

func (r *Router) evidenceCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not allowed"})
		return
	}
	var in CreateEvidenceRequest
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error(), Code: "REQ_001"})
		return
	}
	rec, created, err := r.svc.CreateEvidenceWithIdempotency(req.Header.Get("Idempotency-Key"), service.CreateEvidenceInput(in))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error(), Code: "VAL_001"})
		return
	}
	if !created {
		writeJSON(w, http.StatusOK, CreateEvidenceResponse(rec))
		return
	}
	writeJSON(w, http.StatusCreated, CreateEvidenceResponse(rec))
}

func (r *Router) evidenceGet(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not allowed"})
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/v1/evidence/")
	rec, ok := r.svc.GetEvidence(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, APIError{Error: "not found", Code: "EVID_404"})
		return
	}
	writeJSON(w, http.StatusOK, GetEvidenceResponse(rec))
}

func (r *Router) custodyAppend(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not allowed"})
		return
	}
	var evt CustodyAppendRequest
	if err := json.NewDecoder(req.Body).Decode(&evt); err != nil {
		writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error(), Code: "REQ_001"})
		return
	}
	if err := r.svc.AppendCustody(evt); err != nil {
		writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error(), Code: "VAL_002"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (r *Router) verify(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not allowed"})
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/v1/verify/")
	report := r.svc.VerifyEvidence(id)
	if report.FailureReason == "not_found" {
		writeJSON(w, http.StatusNotFound, VerifyResponse(report))
		return
	}
	writeJSON(w, http.StatusOK, VerifyResponse(report))
}

func (r *Router) search(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not allowed"})
		return
	}
	query := ParseListEvidenceQuery(
		req.URL.Query().Get("page"),
		req.URL.Query().Get("page_size"),
		req.URL.Query().Get("q"),
		req.URL.Query().Get("source"),
		req.URL.Query().Get("type"),
		req.URL.Query().Get("sort_by"),
		req.URL.Query().Get("sort_order"),
	)
	result := r.svc.ListEvidence(query)
	writeJSON(w, http.StatusOK, ListEvidenceResponse{
		Items:    result.Items,
		Page:     result.Page,
		PageSize: result.PageSize,
		Total:    result.Total,
	})
}

func (r *Router) audit(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, AuditResponse{Entries: r.svc.AuditEntries()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
