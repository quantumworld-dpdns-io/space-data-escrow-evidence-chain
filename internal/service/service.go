package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
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
	idemMu   sync.Mutex
	idem     map[string]string
}

func New(e repo.EvidenceRepository, c repo.CustodyRepository, a repo.AuditRepository) *Service {
	return &Service{evidence: e, custody: c, audit: a, idem: map[string]string{}}
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
	for _, existing := range s.evidence.List() {
		if existing.ExternalID == input.ExternalID && existing.Source == input.Source && existing.Type == input.Type {
			return domain.EvidenceRecord{}, errors.New("duplicate_evidence")
		}
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

func (s *Service) CreateEvidenceWithIdempotency(idempotencyKey string, input CreateEvidenceInput) (domain.EvidenceRecord, bool, error) {
	if idempotencyKey == "" {
		rec, err := s.CreateEvidence(input)
		return rec, true, err
	}

	s.idemMu.Lock()
	defer s.idemMu.Unlock()

	if id, ok := s.idem[idempotencyKey]; ok {
		if rec, found := s.evidence.Get(id); found {
			return rec, false, nil
		}
	}

	rec, err := s.CreateEvidence(input)
	if err != nil {
		return domain.EvidenceRecord{}, false, err
	}
	s.idem[idempotencyKey] = rec.ID
	return rec, true, nil
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
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if err := s.custody.Append(event); err != nil {
		return err
	}
	_ = s.audit.Add("custody_appended:" + event.EvidenceID)
	return nil
}

func (s *Service) VerifyEvidence(id string) domain.VerificationReport {
	now := time.Now().UTC().Format(time.RFC3339)
	rec, ok := s.evidence.Get(id)
	if !ok {
		return domain.VerificationReport{EvidenceID: id, FailureReason: "not_found", VerifiedAt: now}
	}

	custody := s.custody.ListByEvidenceID(id)
	if len(custody) == 0 {
		return domain.VerificationReport{
			EvidenceID:       id,
			ChainValid:       false,
			SignatureValid:   true,
			IntegrityValid:   false,
			FailureReason:    "missing_chain_or_custody",
			VerifiedAt:       now,
			CanonicalPayload: crypto.CanonicalizePayload(rec.Payload),
		}
	}

	recomputed := crypto.HashPayload(rec.Payload)
	if recomputed != rec.Hash {
		return domain.VerificationReport{
			EvidenceID:       id,
			ChainValid:       false,
			SignatureValid:   true,
			IntegrityValid:   false,
			FailureReason:    "hash_mismatch",
			VerifiedAt:       now,
			LastCustodyActor: custody[len(custody)-1].Actor,
			CanonicalPayload: crypto.CanonicalizePayload(rec.Payload),
		}
	}

	times := make([]time.Time, 0, len(custody))
	for _, c := range custody {
		if !chain.HasRequiredCustodyFields(c.Actor, c.Action) {
			return domain.VerificationReport{
				EvidenceID:       id,
				ChainValid:       false,
				SignatureValid:   true,
				IntegrityValid:   false,
				FailureReason:    "invalid_custody_event",
				VerifiedAt:       now,
				LastCustodyActor: custody[len(custody)-1].Actor,
				CanonicalPayload: crypto.CanonicalizePayload(rec.Payload),
			}
		}
		times = append(times, c.Timestamp)
	}

	if !chain.HasMonotonicCustodyTimestamps(times) {
		return domain.VerificationReport{
			EvidenceID:       id,
			ChainValid:       false,
			SignatureValid:   true,
			IntegrityValid:   false,
			FailureReason:    "non_monotonic_custody_time",
			VerifiedAt:       now,
			LastCustodyActor: custody[len(custody)-1].Actor,
			CanonicalPayload: crypto.CanonicalizePayload(rec.Payload),
		}
	}

	valid := chain.IsValidChain(rec.Hash, len(custody))
	report := domain.VerificationReport{
		EvidenceID:       id,
		ChainValid:       valid,
		SignatureValid:   true,
		IntegrityValid:   valid,
		VerifiedAt:       now,
		LastCustodyActor: custody[len(custody)-1].Actor,
		CanonicalPayload: crypto.CanonicalizePayload(rec.Payload),
		ComputedHash:     recomputed,
	}
	if !valid {
		report.FailureReason = "missing_chain_or_custody"
	}
	_ = s.audit.Add("evidence_verified:" + id)
	return report
}

func (s *Service) SearchEvidence(q string) []domain.EvidenceRecord {
	return s.ListEvidence(ListEvidenceQuery{Q: q, Page: 1, PageSize: 100, SortBy: "created_at", SortOrder: "desc"}).Items
}

type ListEvidenceResult struct {
	Items    []domain.EvidenceRecord `json:"items"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Total    int                     `json:"total"`
}

func (s *Service) ListEvidence(query ListEvidenceQuery) ListEvidenceResult {
	query.Normalize()
	all := s.evidence.List()
	filtered := make([]domain.EvidenceRecord, 0, len(all))
	for _, r := range all {
		if query.Q != "" && !(r.ExternalID == query.Q || r.Source == query.Q || r.Type == query.Q) {
			continue
		}
		if query.Source != "" && !strings.EqualFold(r.Source, query.Source) {
			continue
		}
		if query.Type != "" && !strings.EqualFold(r.Type, query.Type) {
			continue
		}
		filtered = append(filtered, r)
	}

	sort.Slice(filtered, func(i, j int) bool {
		less := filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		if query.SortBy == "external_id" {
			less = filtered[i].ExternalID < filtered[j].ExternalID
		}
		if query.SortOrder == "asc" {
			return less
		}
		return !less
	})

	total := len(filtered)
	start := (query.Page - 1) * query.PageSize
	if start > total {
		start = total
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return ListEvidenceResult{
		Items:    filtered[start:end],
		Page:     query.Page,
		PageSize: query.PageSize,
		Total:    total,
	}
}

func (s *Service) AuditEntries() []string { return s.audit.List() }
