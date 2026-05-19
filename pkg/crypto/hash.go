package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func HashPayload(payload map[string]string) string {
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+payload[k])
	}
	canonical := strings.Join(parts, "&")
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}
