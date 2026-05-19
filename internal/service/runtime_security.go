package service

import (
	"encoding/json"
	"fmt"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
)

func (s *Service) IngestRuntimeSecurityEvent(evt domain.RuntimeSecurityEvent) error {
	if evt.Source == "" || evt.EventType == "" || evt.Severity == "" || evt.Timestamp == "" {
		return fmt.Errorf("invalid_runtime_event")
	}
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	prefix := "runtime_security/" + evt.Severity + ":"
	_ = s.audit.Add(prefix + string(b))
	s.metrics.Inc("runtime_security.event." + evt.Severity)
	return nil
}
