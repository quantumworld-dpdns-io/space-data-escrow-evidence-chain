package service

import (
	"testing"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
)

func TestSemanticSearchAndEnrichment(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	rec, err := svc.CreateEvidence(CreateEvidenceInput{
		ExternalID: "EXT-SEM",
		Source:     "sat-sem",
		Type:       "imagery",
		Payload:    map[string]string{"desc": "wildfire plume"},
	})
	if err != nil {
		t.Fatalf("create evidence: %v", err)
	}

	items, err := svc.SemanticSearch("wildfire", 10)
	if err != nil {
		t.Fatalf("semantic search: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected semantic search hits")
	}

	job, err := svc.TriggerEnrichment(rec.ID)
	if err != nil {
		t.Fatalf("trigger enrichment: %v", err)
	}
	if job.Status == "" {
		t.Fatalf("expected job status: %+v", job)
	}
	stored, ok := svc.GetEnrichmentJob(job.ID)
	if !ok {
		t.Fatal("expected job status lookup success")
	}
	if stored.Status != job.Status {
		t.Fatalf("job mismatch: %+v vs %+v", job, stored)
	}
}
