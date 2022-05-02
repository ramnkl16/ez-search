package ezsearch

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

type SearchResponseModel struct {
	ResultRow []map[string]interface{} `json:"resultRow"`
	Fields    []string                 `json:"fields"`
	Facets    map[string][]EzTermFacet `json:"facetResult"`
	Total     uint64                   `json:"total"`
	Took      time.Duration            `json:"took"`
	Status    bleve.SearchStatus       `json:"status"`
}

type EzTermFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}
