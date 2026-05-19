package repo

import "github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"

type EvidenceRepository interface {
	Create(record domain.EvidenceRecord) error
	Get(id string) (domain.EvidenceRecord, bool)
	List() []domain.EvidenceRecord
}

type CustodyRepository interface {
	Append(event domain.CustodyEvent) error
	ListByEvidenceID(id string) []domain.CustodyEvent
}

type AuditRepository interface {
	Add(entry string) error
	List() []string
}
