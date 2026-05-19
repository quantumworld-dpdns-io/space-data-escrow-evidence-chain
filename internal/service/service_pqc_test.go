package service

import (
	"testing"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
)

func TestAddAttestationDualSignPlaceholder(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	rec, err := svc.CreateEvidence(CreateEvidenceInput{ExternalID: "EXT-PQC", Source: "sat-pqc", Type: "img", Payload: map[string]string{"k": "v"}})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	att, err := svc.AddAttestation(domain.Attestation{
		EvidenceID: rec.ID,
		Signer:     "ops",
		Algorithm:  "ed25519",
		DualSign:   true,
	})
	if err != nil {
		t.Fatalf("attest: %v", err)
	}
	if att.ClassicalSignature == "" || att.PQCSignature == "" || att.PQCAlgorithm == "" {
		t.Fatalf("expected dual signatures and pqc algorithm: %+v", att)
	}
}
