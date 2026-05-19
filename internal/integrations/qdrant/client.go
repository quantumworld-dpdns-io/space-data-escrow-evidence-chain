package qdrant

import (
	"context"
	"strings"
)

type SearchResult struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
}

type Client interface {
	Upsert(ctx context.Context, id string, payload map[string]string) error
	Search(ctx context.Context, query string, limit int) ([]SearchResult, error)
}

type MemoryClient struct {
	data map[string]map[string]string
}

func NewMemoryClient() *MemoryClient {
	return &MemoryClient{data: map[string]map[string]string{}}
}

func (c *MemoryClient) Upsert(_ context.Context, id string, payload map[string]string) error {
	c.data[id] = payload
	return nil
}

func (c *MemoryClient) Search(_ context.Context, query string, limit int) ([]SearchResult, error) {
	out := make([]SearchResult, 0)
	q := strings.ToLower(query)
	for id, p := range c.data {
		hit := false
		for _, v := range p {
			if strings.Contains(strings.ToLower(v), q) {
				hit = true
				break
			}
		}
		if hit {
			out = append(out, SearchResult{ID: id, Score: 1.0})
		}
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}
