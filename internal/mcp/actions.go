package mcp

import (
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

type Actions struct {
	svc *service.Service
}

func NewActions(svc *service.Service) *Actions {
	return &Actions{svc: svc}
}

func (a *Actions) EvidenceIngest(in service.CreateEvidenceInput) (domain.EvidenceRecord, error) {
	rec, _, err := a.svc.CreateEvidenceWithIdempotency("", in)
	return rec, err
}

func (a *Actions) ChainVerify(id string) domain.VerificationReport {
	return a.svc.VerifyEvidence(id)
}

func (a *Actions) SemanticSearch(q string, limit int) ([]domain.EvidenceRecord, error) {
	return a.svc.SemanticSearch(q, limit)
}

func (a *Actions) AuditQuery() []string {
	return a.svc.AuditEntries()
}
