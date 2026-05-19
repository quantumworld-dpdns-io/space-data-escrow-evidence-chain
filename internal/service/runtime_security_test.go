package service

import (
	"testing"
	"time"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
)

func TestIngestRuntimeSecurityEvent(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	err := svc.IngestRuntimeSecurityEvent(domain.RuntimeSecurityEvent{
		Source:    "tetragon",
		EventType: "suspicious_exec",
		Severity:  "high",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := svc.AuditEntries()
	if len(entries) == 0 {
		t.Fatal("expected runtime audit entry")
	}
}

func TestIngestRuntimeSecurityEventValidation(t *testing.T) {
	svc := New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	err := svc.IngestRuntimeSecurityEvent(domain.RuntimeSecurityEvent{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
