package cli

func (c *Client) DoEvidenceCreate(externalID, source, typ string, payload map[string]string) ([]byte, int, error) {
	return c.do("POST", "/v1/evidence", map[string]any{
		"external_id": externalID,
		"source":      source,
		"type":        typ,
		"payload":     payload,
	})
}

func (c *Client) DoEvidenceVerify(id string) ([]byte, int, error) {
	return c.do("POST", "/v1/verify/"+id, nil)
}

func (c *Client) DoAuditQuery() ([]byte, int, error) {
	return c.do("GET", "/v1/audit", nil)
}

func (c *Client) DoEnrichTrigger(id string) ([]byte, int, error) {
	return c.do("POST", "/v1/enrich", map[string]any{"evidence_id": id})
}
