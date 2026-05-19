package service

import (
	"testing"
	"time"

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

	err = svc.AppendCustody(domain.CustodyEvent{EvidenceID: rec.ID, Actor: "operator-1", Action: "ingest"})
	if err != nil {
		t.Fatalf("append custody: %v", err)
	}

	report := svc.VerifyEvidence(rec.ID)
	if !report.ChainValid || !report.IntegrityValid {
		t.Fatalf("expected valid chain report: %+v", report)
	}
	if report.CanonicalPayload == "" || report.ComputedHash == "" {
		t.Fatalf("expected canonical payload and computed hash: %+v", report)
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

func TestVerifyEvidenceMissingCustody(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	rec, _ := svc.CreateEvidence(CreateEvidenceInput{ExternalID: "EXT-2", Source: "sat-b", Type: "telemetry", Payload: map[string]string{"k": "v"}})
	report := svc.VerifyEvidence(rec.ID)
	if report.FailureReason != "missing_chain_or_custody" {
		t.Fatalf("expected missing_chain_or_custody, got %+v", report)
	}
}

func TestVerifyEvidenceHashMismatch(t *testing.T) {
	e := memory.NewEvidenceRepo()
	c := memory.NewCustodyRepo()
	a := memory.NewAuditRepo()
	svc := New(e, c, a)

	rec, _ := svc.CreateEvidence(CreateEvidenceInput{ExternalID: "EXT-3", Source: "sat-c", Type: "imagery", Payload: map[string]string{"f": "1"}})
	_ = c.Append(domain.CustodyEvent{EvidenceID: rec.ID, Actor: "op", Action: "ingest", Timestamp: time.Now().UTC()})

	stored, _ := e.Get(rec.ID)
	stored.Hash = "tampered"
	_ = e.Create(stored)

	report := svc.VerifyEvidence(rec.ID)
	if report.FailureReason != "hash_mismatch" {
		t.Fatalf("expected hash_mismatch, got %+v", report)
	}
}

func TestVerifyEvidenceNonMonotonicCustody(t *testing.T) {
	e := memory.NewEvidenceRepo()
	c := memory.NewCustodyRepo()
	a := memory.NewAuditRepo()
	svc := New(e, c, a)

	rec, _ := svc.CreateEvidence(CreateEvidenceInput{ExternalID: "EXT-4", Source: "sat-d", Type: "imagery", Payload: map[string]string{"f": "1"}})
	now := time.Now().UTC()
	_ = c.Append(domain.CustodyEvent{EvidenceID: rec.ID, Actor: "op1", Action: "ingest", Timestamp: now})
	_ = c.Append(domain.CustodyEvent{EvidenceID: rec.ID, Actor: "op2", Action: "transfer", Timestamp: now.Add(-time.Minute)})

	report := svc.VerifyEvidence(rec.ID)
	if report.FailureReason != "non_monotonic_custody_time" {
		t.Fatalf("expected non_monotonic_custody_time, got %+v", report)
	}
}
