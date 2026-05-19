package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

type Router struct {
	svc    *service.Service
	apiKey string
	mux    *http.ServeMux
}

func NewRouter(svc *service.Service, apiKey string) *Router {
	r := &Router{svc: svc, apiKey: apiKey, mux: http.NewServeMux()}
	r.routes()
	return r
}

func (r *Router) Handler() http.Handler { return r.mux }

func (r *Router) routes() {
	r.mux.HandleFunc("/healthz", r.health)
	r.mux.HandleFunc("/readyz", r.ready)
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
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		next(w, req)
	}
}

func (r *Router) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}
func (r *Router) ready(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status": "ready"}) }

func (r *Router) evidenceCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}); return }
	var in service.CreateEvidenceInput
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); return }
	rec, err := r.svc.CreateEvidence(in)
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); return }
	writeJSON(w, http.StatusCreated, rec)
}

func (r *Router) evidenceGet(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}); return }
	id := strings.TrimPrefix(req.URL.Path, "/v1/evidence/")
	rec, ok := r.svc.GetEvidence(id)
	if !ok { writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"}); return }
	writeJSON(w, http.StatusOK, rec)
}

func (r *Router) custodyAppend(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}); return }
	var evt domain.CustodyEvent
	if err := json.NewDecoder(req.Body).Decode(&evt); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); return }
	if err := r.svc.AppendCustody(evt); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); return }
	w.WriteHeader(http.StatusAccepted)
}

func (r *Router) verify(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}); return }
	id := strings.TrimPrefix(req.URL.Path, "/v1/verify/")
	report := r.svc.VerifyEvidence(id)
	if report.FailureReason == "not_found" { writeJSON(w, http.StatusNotFound, report); return }
	writeJSON(w, http.StatusOK, report)
}

func (r *Router) search(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}); return }
	writeJSON(w, http.StatusOK, r.svc.SearchEvidence(req.URL.Query().Get("q")))
}

func (r *Router) audit(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}); return }
	writeJSON(w, http.StatusOK, map[string]any{"entries": r.svc.AuditEntries()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
