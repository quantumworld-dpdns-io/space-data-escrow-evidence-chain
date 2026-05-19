package mcp

import (
	"testing"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

func TestMCPActionsFlow(t *testing.T) {
	svc := service.New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	a := NewActions(svc)

	rec, err := a.EvidenceIngest(service.CreateEvidenceInput{
		ExternalID: "EXT-MCP",
		Source:     "sat-mcp",
		Type:       "imagery",
		Payload:    map[string]string{"desc": "launch"},
	})
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}

	report := a.ChainVerify(rec.ID)
	if report.EvidenceID == "" {
		t.Fatalf("verify empty report: %+v", report)
	}

	items, err := a.SemanticSearch("launch", 5)
	if err != nil {
		t.Fatalf("semantic search: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected semantic search results")
	}

	audit := a.AuditQuery()
	if len(audit) == 0 {
		t.Fatal("expected audit entries")
	}
}
