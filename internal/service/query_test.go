package service

import "testing"

func TestListEvidenceQueryNormalize(t *testing.T) {
	q := ListEvidenceQuery{}
	q.Normalize()
	if q.Page != 1 || q.PageSize != 20 || q.SortBy != "created_at" || q.SortOrder != "desc" {
		t.Fatalf("unexpected normalized query: %+v", q)
	}
}
