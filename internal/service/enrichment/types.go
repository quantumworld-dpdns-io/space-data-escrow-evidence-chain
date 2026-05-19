package enrichment

import "time"

type JobStatus string

const (
	JobPending JobStatus = "pending"
	JobDone    JobStatus = "done"
	JobFailed  JobStatus = "failed"
)

type Job struct {
	ID         string            `json:"id"`
	EvidenceID string            `json:"evidence_id"`
	Status     JobStatus         `json:"status"`
	Output     map[string]string `json:"output,omitempty"`
	Error      string            `json:"error,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
