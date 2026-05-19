package crypto

import "testing"

func TestCanonicalizePayloadDeterministic(t *testing.T) {
	a := map[string]string{"b": "2", "a": "1"}
	b := map[string]string{"a": "1", "b": "2"}
	ca := CanonicalizePayload(a)
	cb := CanonicalizePayload(b)
	if ca != cb {
		t.Fatalf("expected same canonical payload, got %q vs %q", ca, cb)
	}
	if HashPayload(a) != HashPayload(b) {
		t.Fatalf("expected same hash for equivalent payload maps")
	}
}
