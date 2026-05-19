package crypto

type PQCSigner interface {
	Algorithm() string
	Sign(message []byte) (string, error)
	Verify(message []byte, signature string) bool
}

type DilithiumSigner struct{}

func (s DilithiumSigner) Algorithm() string { return "dilithium-placeholder" }
func (s DilithiumSigner) Sign(_ []byte) (string, error) { return "dilithium-signature-placeholder", nil }
func (s DilithiumSigner) Verify(_ []byte, signature string) bool {
	return signature == "dilithium-signature-placeholder"
}

type KyberKEM struct{}

func (k KyberKEM) Algorithm() string { return "kyber-placeholder" }
func (k KyberKEM) Encapsulate() string { return "kyber-ciphertext-placeholder" }
func (k KyberKEM) Decapsulate(_ string) string { return "kyber-shared-secret-placeholder" }
