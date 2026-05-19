package crypto

type HashEngine interface {
	Name() string
	Hash(payload map[string]string) string
}

type Signer interface {
	Algorithm() string
	Sign(message []byte) (string, error)
	Verify(message []byte, signature string) bool
}
