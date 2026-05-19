package telemetry

import "testing"

func TestRegistryCountersAndP95(t *testing.T) {
	r := NewRegistry()
	r.Inc("a")
	r.Inc("a")
	r.AddLatency("l", 10)
	r.AddLatency("l", 50)
	r.AddLatency("l", 30)
	if r.SnapshotCounters()["a"] != 2 {
		t.Fatalf("expected counter=2 got %d", r.SnapshotCounters()["a"])
	}
	if p := r.LatencyP95("l"); p <= 0 {
		t.Fatalf("expected p95 > 0 got %d", p)
	}
}
