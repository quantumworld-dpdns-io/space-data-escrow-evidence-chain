package cli

import "testing"

func TestClientCtor(t *testing.T) {
	c := New("http://localhost:8080", "k")
	if c.BaseURL == "" || c.APIKey == "" || c.HTTP == nil {
		t.Fatal("invalid client initialization")
	}
}
