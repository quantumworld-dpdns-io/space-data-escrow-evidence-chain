package api

import "testing"

func TestParseListEvidenceQueryDefaults(t *testing.T) {
	q := ParseListEvidenceQuery("", "", "", "", "", "", "")
	if q.Page != 1 || q.PageSize != 20 {
		t.Fatalf("unexpected defaults: %+v", q)
	}
	if q.SortBy != "created_at" || q.SortOrder != "desc" {
		t.Fatalf("unexpected sort defaults: %+v", q)
	}
}
