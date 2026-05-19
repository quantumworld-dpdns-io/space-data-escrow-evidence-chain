package service

import (
	"testing"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
)

func TestAttestationAndProofBundle(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	rec, err := svc.CreateEvidence(CreateEvidenceInput{ExternalID: "EXT-PROOF", Source: "sat-p", Type: "imagery", Payload: map[string]string{"k": "v"}})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	_, _ = svc.AppendCustody(domain.CustodyEvent{EvidenceID: rec.ID, Actor: "op", Action: "ingest"}), nil
	_, err = svc.AddAttestation(domain.Attestation{EvidenceID: rec.ID, Signer: "s1", Signature: "sig", Algorithm: "ed25519"})
	if err != nil {
		t.Fatalf("attestation: %v", err)
	}
	bundle, err := svc.ExportProofBundle(rec.ID)
	if err != nil {
		t.Fatalf("proof: %v", err)
	}
	if bundle.Evidence.ID == "" || len(bundle.Attestations) == 0 {
		t.Fatalf("unexpected bundle: %+v", bundle)
	}
}

func TestChainTimelineNotFound(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	_, err := svc.ChainTimeline("missing")
	if err == nil {
		t.Fatal("expected error")
	}
}
