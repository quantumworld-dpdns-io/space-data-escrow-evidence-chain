package telemetry

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPMetricsMiddleware(t *testing.T) {
	reg := NewRegistry()
	mw := HTTPMetrics(reg)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	c := reg.SnapshotCounters()
	if c["http.requests.total"] == 0 {
		t.Fatalf("expected request counter increment: %+v", c)
	}
}
