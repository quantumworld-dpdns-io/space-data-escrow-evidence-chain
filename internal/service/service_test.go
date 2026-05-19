package service

import (
	"testing"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
)

func TestEvidenceLifecycle(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())

	rec, err := svc.CreateEvidence(CreateEvidenceInput{
		ExternalID: "EXT-1",
		Source:     "satellite-a",
		Type:       "imagery",
		Payload:    map[string]string{"file": "a.tif"},
	})
	if err != nil {
		t.Fatalf("create evidence: %v", err)
	}

	if rec.Hash == "" {
		t.Fatal("expected hash")
	}

	err = svc.AppendCustody(domain.CustodyEvent{EvidenceID: rec.ID, Actor: "operator-1", Action: "ingest"})
	if err != nil {
		t.Fatalf("append custody: %v", err)
	}

	report := svc.VerifyEvidence(rec.ID)
	if !report.ChainValid || !report.IntegrityValid {
		t.Fatalf("expected valid chain report: %+v", report)
	}
}

func TestCreateEvidenceValidation(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	_, err := svc.CreateEvidence(CreateEvidenceInput{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestVerifyEvidenceNotFound(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	report := svc.VerifyEvidence("missing")
	if report.FailureReason != "not_found" {
		t.Fatalf("expected not_found, got %+v", report)
	}
}
