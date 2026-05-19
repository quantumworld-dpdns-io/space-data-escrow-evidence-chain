package telemetry

import (
	"context"
	"fmt"
	"time"
)

type spanKey struct{}

type Span struct {
	Name      string
	StartedAt time.Time
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	sp := Span{Name: name, StartedAt: time.Now().UTC()}
	return context.WithValue(ctx, spanKey{}, sp), sp
}

func EndSpan(reg *Registry, sp Span) {
	ms := time.Since(sp.StartedAt).Milliseconds()
	reg.AddLatency("span."+sp.Name, ms)
	reg.Inc("span.count")
}

func RequestIDFromTime() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
