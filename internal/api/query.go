package api

import (
	"strconv"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

func ParseListEvidenceQuery(pageRaw, pageSizeRaw, q, source, typ, sortBy, sortOrder string) service.ListEvidenceQuery {
	page, _ := strconv.Atoi(pageRaw)
	pageSize, _ := strconv.Atoi(pageSizeRaw)
	query := service.ListEvidenceQuery{
		Q:         q,
		Source:    source,
		Type:      typ,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      page,
		PageSize:  pageSize,
	}
	query.Normalize()
	return query
}
