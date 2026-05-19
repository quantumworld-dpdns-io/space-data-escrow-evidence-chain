package chain

import "time"

type Link struct {
	EvidenceID string
	Hash       string
	Actor      string
	Action     string
	Timestamp  time.Time
}
