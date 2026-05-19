package telemetry

import (
	"sort"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	counters map[string]int64
	latency  map[string][]int64
}

func NewRegistry() *Registry {
	return &Registry{counters: map[string]int64{}, latency: map[string][]int64{}}
}

func (r *Registry) Inc(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[key]++
}

func (r *Registry) AddLatency(key string, ms int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latency[key] = append(r.latency[key], ms)
}

func (r *Registry) SnapshotCounters() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]int64, len(r.counters))
	for k, v := range r.counters {
		out[k] = v
	}
	return out
}

func (r *Registry) LatencyP95(key string) int64 {
	r.mu.RLock()
	vals := append([]int64(nil), r.latency[key]...)
	r.mu.RUnlock()
	if len(vals) == 0 {
		return 0
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	idx := int(float64(len(vals)-1) * 0.95)
	return vals[idx]
}
