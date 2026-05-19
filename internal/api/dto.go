package api

import "github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"

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
