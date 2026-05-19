package domain

import "time"

type EvidenceRecord struct {
	ID         string            `json:"id"`
	ExternalID string            `json:"external_id"`
	Source     string            `json:"source"`
	Type       string            `json:"type"`
	Payload    map[string]string `json:"payload"`
	Hash       string            `json:"hash"`
	CreatedAt  time.Time         `json:"created_at"`
}

type CustodyEvent struct {
	EvidenceID string    `json:"evidence_id"`
	Actor      string    `json:"actor"`
	Action     string    `json:"action"`
	Timestamp  time.Time `json:"timestamp"`
	Note       string    `json:"note"`
}

type Attestation struct {
	EvidenceID string    `json:"evidence_id"`
	Signer     string    `json:"signer"`
	Signature  string    `json:"signature"`
	Algorithm  string    `json:"algorithm"`
	Timestamp  time.Time `json:"timestamp"`
}

type VerificationReport struct {
	EvidenceID        string `json:"evidence_id"`
	ChainValid        bool   `json:"chain_valid"`
	SignatureValid    bool   `json:"signature_valid"`
	IntegrityValid    bool   `json:"integrity_valid"`
	FailureReason     string `json:"failure_reason,omitempty"`
	VerifiedAt        string `json:"verified_at"`
	LastCustodyActor  string `json:"last_custody_actor,omitempty"`
	CanonicalPayload  string `json:"canonical_payload,omitempty"`
	ComputedHash      string `json:"computed_hash,omitempty"`
}
