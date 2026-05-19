package ollama

import "context"

type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type MemoryClient struct{}

func NewMemoryClient() *MemoryClient { return &MemoryClient{} }

func (c *MemoryClient) Generate(_ context.Context, prompt string) (string, error) {
	return "enriched:" + prompt, nil
}
