package memory

import (
	"sync"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
)

type EvidenceRepo struct {
	mu   sync.RWMutex
	data map[string]domain.EvidenceRecord
}

func NewEvidenceRepo() *EvidenceRepo {
	return &EvidenceRepo{data: map[string]domain.EvidenceRecord{}}
}

func (r *EvidenceRepo) Create(record domain.EvidenceRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[record.ID] = record
	return nil
}

func (r *EvidenceRepo) Get(id string) (domain.EvidenceRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.data[id]
	return rec, ok
}

func (r *EvidenceRepo) List() []domain.EvidenceRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.EvidenceRecord, 0, len(r.data))
	for _, v := range r.data {
		out = append(out, v)
	}
	return out
}

type CustodyRepo struct {
	mu   sync.RWMutex
	data map[string][]domain.CustodyEvent
}

func NewCustodyRepo() *CustodyRepo {
	return &CustodyRepo{data: map[string][]domain.CustodyEvent{}}
}

func (r *CustodyRepo) Append(event domain.CustodyEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[event.EvidenceID] = append(r.data[event.EvidenceID], event)
	return nil
}

func (r *CustodyRepo) ListByEvidenceID(id string) []domain.CustodyEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make([]domain.CustodyEvent, len(r.data[id]))
	copy(cp, r.data[id])
	return cp
}

type AuditRepo struct {
	mu      sync.RWMutex
	entries []string
}

func NewAuditRepo() *AuditRepo {
	return &AuditRepo{entries: []string{}}
}

func (r *AuditRepo) Add(entry string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry)
	return nil
}

func (r *AuditRepo) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make([]string, len(r.entries))
	copy(cp, r.entries)
	return cp
}
