package crypto

type NoopSigner struct{}

func (s NoopSigner) Algorithm() string { return "noop" }

func (s NoopSigner) Sign(_ []byte) (string, error) { return "noop-signature", nil }

func (s NoopSigner) Verify(_ []byte, signature string) bool { return signature == "noop-signature" }
