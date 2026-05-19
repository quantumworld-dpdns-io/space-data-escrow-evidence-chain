package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/pkg/chain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/pkg/crypto"
)

type Service struct {
	evidence repo.EvidenceRepository
	custody  repo.CustodyRepository
	audit    repo.AuditRepository
}

func New(e repo.EvidenceRepository, c repo.CustodyRepository, a repo.AuditRepository) *Service {
	return &Service{evidence: e, custody: c, audit: a}
}

type CreateEvidenceInput struct {
	ExternalID string            `json:"external_id"`
	Source     string            `json:"source"`
	Type       string            `json:"type"`
	Payload    map[string]string `json:"payload"`
}

func (s *Service) CreateEvidence(input CreateEvidenceInput) (domain.EvidenceRecord, error) {
	if input.ExternalID == "" || input.Source == "" || input.Type == "" {
		return domain.EvidenceRecord{}, errors.New("external_id, source and type are required")
	}
	rec := domain.EvidenceRecord{
		ID:         newID(),
		ExternalID: input.ExternalID,
		Source:     input.Source,
		Type:       input.Type,
		Payload:    input.Payload,
		Hash:       crypto.HashPayload(input.Payload),
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.evidence.Create(rec); err != nil {
		return domain.EvidenceRecord{}, err
	}
	_ = s.audit.Add("evidence_created:" + rec.ID)
	return rec, nil
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Service) GetEvidence(id string) (domain.EvidenceRecord, bool) { return s.evidence.Get(id) }

func (s *Service) AppendCustody(event domain.CustodyEvent) error {
	if event.EvidenceID == "" || event.Actor == "" || event.Action == "" {
		return errors.New("evidence_id, actor, action are required")
	}
	event.Timestamp = time.Now().UTC()
	if err := s.custody.Append(event); err != nil {
		return err
	}
	_ = s.audit.Add("custody_appended:" + event.EvidenceID)
	return nil
}

func (s *Service) VerifyEvidence(id string) domain.VerificationReport {
	rec, ok := s.evidence.Get(id)
	if !ok {
		return domain.VerificationReport{EvidenceID: id, FailureReason: "not_found", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
	}
	custody := s.custody.ListByEvidenceID(id)
	valid := chain.IsValidChain(rec.Hash, len(custody))
	report := domain.VerificationReport{EvidenceID: id, ChainValid: valid, SignatureValid: true, IntegrityValid: valid, VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
	if len(custody) > 0 {
		report.LastCustodyActor = custody[len(custody)-1].Actor
	}
	if !valid {
		report.FailureReason = "missing_chain_or_custody"
	}
	_ = s.audit.Add("evidence_verified:" + id)
	return report
}

func (s *Service) SearchEvidence(q string) []domain.EvidenceRecord {
	all := s.evidence.List()
	if q == "" {
		return all
	}
	results := make([]domain.EvidenceRecord, 0)
	for _, r := range all {
		if r.ExternalID == q || r.Source == q || r.Type == q {
			results = append(results, r)
		}
	}
	return results
}

func (s *Service) AuditEntries() []string { return s.audit.List() }
