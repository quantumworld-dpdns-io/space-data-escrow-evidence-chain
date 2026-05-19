package api

import (
	"time"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
)

type CreateEvidenceRequest struct {
	ExternalID string            `json:"external_id"`
	Source     string            `json:"source"`
	Type       string            `json:"type"`
	Payload    map[string]string `json:"payload"`
}

type CreateEvidenceResponse = domain.EvidenceRecord

type GetEvidenceResponse = domain.EvidenceRecord

type CustodyAppendRequest = domain.CustodyEvent

type VerifyResponse = domain.VerificationReport

type AuditResponse struct {
	Entries []string `json:"entries"`
}

type ListEvidenceResponse struct {
	Items    []domain.EvidenceRecord `json:"items"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Total    int                     `json:"total"`
}

type TriggerEnrichmentRequest struct {
	EvidenceID string `json:"evidence_id"`
}

type EnrichmentJobResponse struct {
	ID         string            `json:"id"`
	EvidenceID string            `json:"evidence_id"`
	Status     string            `json:"status"`
	Output     map[string]string `json:"output,omitempty"`
	Error      string            `json:"error,omitempty"`
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
}

type AttestationRequest struct {
	EvidenceID string `json:"evidence_id"`
	Signer     string `json:"signer"`
	Signature  string `json:"signature"`
	Algorithm  string `json:"algorithm"`
	Timestamp  string `json:"timestamp,omitempty"`
}

func (r AttestationRequest) ToDomain() domain.Attestation {
	var ts time.Time
	if r.Timestamp != "" {
		parsed, err := time.Parse(time.RFC3339, r.Timestamp)
		if err == nil {
			ts = parsed
		}
	}
	return domain.Attestation{
		EvidenceID: r.EvidenceID,
		Signer:     r.Signer,
		Signature:  r.Signature,
		Algorithm:  r.Algorithm,
		Timestamp:  ts,
	}
}
