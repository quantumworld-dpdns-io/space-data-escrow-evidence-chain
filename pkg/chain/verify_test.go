package chain

import (
	"testing"
	"time"
)

func TestHasMonotonicCustodyTimestamps(t *testing.T) {
	now := time.Now()
	if !HasMonotonicCustodyTimestamps([]time.Time{now, now.Add(time.Second), now.Add(2 * time.Second)}) {
		t.Fatal("expected monotonic timestamps to pass")
	}
	if HasMonotonicCustodyTimestamps([]time.Time{now.Add(time.Second), now}) {
		t.Fatal("expected non-monotonic timestamps to fail")
	}
}
