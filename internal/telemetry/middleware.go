package telemetry

import (
	"net/http"
	"time"
)

func HTTPMetrics(reg *Registry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx, sp := StartSpan(r.Context(), "http.request")
			next.ServeHTTP(w, r.WithContext(ctx))
			EndSpan(reg, sp)
			reg.Inc("http.requests.total")
			reg.AddLatency("http.requests.latency_ms", time.Since(start).Milliseconds())
			reg.Inc("http.route." + r.URL.Path)
		})
	}
}
