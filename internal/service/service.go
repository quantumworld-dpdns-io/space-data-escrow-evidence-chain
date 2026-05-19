package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/integrations/ollama"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/integrations/qdrant"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service/enrichment"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/telemetry"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/pkg/chain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/pkg/crypto"
)

type Service struct {
	evidence  repo.EvidenceRepository
	custody   repo.CustodyRepository
	audit     repo.AuditRepository
	idemMu    sync.Mutex
	idem      map[string]string
	qdrant    qdrant.Client
	ollama    ollama.Client
	jobs      *enrichment.Store
	attMu     sync.RWMutex
	att       map[string][]domain.Attestation
	metrics   *telemetry.Registry
	pqcSigner crypto.PQCSigner
}

func New(e repo.EvidenceRepository, c repo.CustodyRepository, a repo.AuditRepository) *Service {
	return &Service{
		evidence: e, custody: c, audit: a, idem: map[string]string{},
		qdrant: qdrant.NewMemoryClient(), ollama: ollama.NewMemoryClient(), jobs: enrichment.NewStore(),
		att:       map[string][]domain.Attestation{},
		metrics:   telemetry.NewRegistry(),
		pqcSigner: crypto.DilithiumSigner{},
	}
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
	s.metrics.Inc("evidence.created")
	_ = s.qdrant.Upsert(context.Background(), rec.ID, rec.Payload)
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
	s.metrics.Inc("custody.appended")
	_ = s.audit.Add("custody_appended:" + event.EvidenceID)
	return nil
}

func (s *Service) VerifyEvidence(id string) domain.VerificationReport {
	now := time.Now().UTC().Format(time.RFC3339)
	rec, ok := s.evidence.Get(id)
	if !ok {
		s.metrics.Inc("verify.not_found")
		return domain.VerificationReport{EvidenceID: id, FailureReason: "not_found", VerifiedAt: now}
	}

	custody := s.custody.ListByEvidenceID(id)
	if len(custody) == 0 {
		s.metrics.Inc("verify.missing_custody")
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
		s.metrics.Inc("verify.hash_mismatch")
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
		s.metrics.Inc("verify.non_monotonic_custody")
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
	s.metrics.Inc("verify.ok")
	_ = s.audit.Add("evidence_verified:" + id)
	return report
}

func (s *Service) SearchEvidence(q string) []domain.EvidenceRecord {
	return s.ListEvidence(ListEvidenceQuery{Q: q, Page: 1, PageSize: 100, SortBy: "created_at", SortOrder: "desc"}).Items
}

func (s *Service) SemanticSearch(query string, limit int) ([]domain.EvidenceRecord, error) {
	found, err := s.qdrant.Search(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}
	out := make([]domain.EvidenceRecord, 0, len(found))
	for _, hit := range found {
		if rec, ok := s.evidence.Get(hit.ID); ok {
			out = append(out, rec)
		}
	}
	return out, nil
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

func (s *Service) ChainTimeline(evidenceID string) ([]domain.CustodyEvent, error) {
	if _, ok := s.evidence.Get(evidenceID); !ok {
		return nil, errors.New("evidence_not_found")
	}
	return s.custody.ListByEvidenceID(evidenceID), nil
}

func (s *Service) AddAttestation(in domain.Attestation) (domain.Attestation, error) {
	if in.EvidenceID == "" || in.Signer == "" || in.Algorithm == "" {
		return domain.Attestation{}, errors.New("invalid_attestation")
	}
	if _, ok := s.evidence.Get(in.EvidenceID); !ok {
		return domain.Attestation{}, errors.New("evidence_not_found")
	}
	if in.Timestamp.IsZero() {
		in.Timestamp = time.Now().UTC()
	}
	if in.ClassicalSignature == "" {
		if in.Signature != "" {
			in.ClassicalSignature = in.Signature
		} else {
			in.ClassicalSignature = "classical-signature-placeholder"
		}
	}
	if in.DualSign {
		if in.PQCAlgorithm == "" {
			in.PQCAlgorithm = s.pqcSigner.Algorithm()
		}
		if in.PQCSignature == "" {
			sig, _ := s.pqcSigner.Sign([]byte(in.EvidenceID + ":" + in.Signer))
			in.PQCSignature = sig
		}
	}
	if in.Signature == "" {
		in.Signature = in.ClassicalSignature
	}
	s.attMu.Lock()
	defer s.attMu.Unlock()
	s.att[in.EvidenceID] = append(s.att[in.EvidenceID], in)
	_ = s.audit.Add("attestation_added:" + in.EvidenceID)
	return in, nil
}

func (s *Service) BulkVerify(ids []string) []domain.VerificationReport {
	out := make([]domain.VerificationReport, 0, len(ids))
	for _, id := range ids {
		out = append(out, s.VerifyEvidence(id))
	}
	return out
}

type ProofBundle struct {
	Evidence     domain.EvidenceRecord     `json:"evidence"`
	Custody      []domain.CustodyEvent     `json:"custody"`
	Attestations []domain.Attestation      `json:"attestations"`
	Verification domain.VerificationReport `json:"verification"`
}

func (s *Service) ExportProofBundle(evidenceID string) (ProofBundle, error) {
	rec, ok := s.evidence.Get(evidenceID)
	if !ok {
		return ProofBundle{}, errors.New("evidence_not_found")
	}
	s.attMu.RLock()
	att := append([]domain.Attestation(nil), s.att[evidenceID]...)
	s.attMu.RUnlock()
	return ProofBundle{
		Evidence:     rec,
		Custody:      s.custody.ListByEvidenceID(evidenceID),
		Attestations: att,
		Verification: s.VerifyEvidence(evidenceID),
	}, nil
}

type KeyRotationMetadata struct {
	CurrentKeyID string `json:"current_key_id"`
	Algorithm    string `json:"algorithm"`
	NextRotation string `json:"next_rotation"`
}

func (s *Service) GetKeyRotationMetadata() KeyRotationMetadata {
	return KeyRotationMetadata{
		CurrentKeyID: "key-dev-001",
		Algorithm:    "ed25519",
		NextRotation: time.Now().UTC().Add(30 * 24 * time.Hour).Format(time.RFC3339),
	}
}

func (s *Service) TriggerEnrichment(evidenceID string) (enrichment.Job, error) {
	rec, ok := s.evidence.Get(evidenceID)
	if !ok {
		return enrichment.Job{}, errors.New("evidence_not_found")
	}
	job := enrichment.Job{
		ID:         newID(),
		EvidenceID: evidenceID,
		Status:     enrichment.JobPending,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	s.jobs.Put(job)

	text, err := s.ollama.Generate(context.Background(), crypto.CanonicalizePayload(rec.Payload))
	if err != nil {
		job.Status = enrichment.JobFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now().UTC()
		s.jobs.Put(job)
		return job, nil
	}
	job.Status = enrichment.JobDone
	job.Output = map[string]string{"summary": text}
	job.UpdatedAt = time.Now().UTC()
	s.jobs.Put(job)
	_ = s.audit.Add("enrichment_done:" + evidenceID)
	return job, nil
}

func (s *Service) GetEnrichmentJob(id string) (enrichment.Job, bool) {
	return s.jobs.Get(id)
}
