package domain

type RuntimeSecurityEvent struct {
	Source    string            `json:"source"`
	EventType string            `json:"event_type"`
	Severity  string            `json:"severity"`
	Actor     string            `json:"actor,omitempty"`
	Resource  string            `json:"resource,omitempty"`
	Timestamp string            `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
}
