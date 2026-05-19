package service

import "strings"

type ListEvidenceQuery struct {
	Q         string
	Source    string
	Type      string
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

func (q *ListEvidenceQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	q.SortBy = strings.ToLower(q.SortBy)
	q.SortOrder = strings.ToLower(q.SortOrder)
}
