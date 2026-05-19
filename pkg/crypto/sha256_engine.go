package crypto

type SHA256Engine struct{}

func (e SHA256Engine) Name() string { return "sha256" }

func (e SHA256Engine) Hash(payload map[string]string) string {
	return HashPayload(payload)
}
