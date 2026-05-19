package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPayload(payload map[string]string) string {
	canonical := CanonicalizePayload(payload)
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}
