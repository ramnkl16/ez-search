package ezsearch

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
)

func buildAggProdIndexMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	dtFieldMapping := bleve.NewDateTimeFieldMapping()

	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name
	//keywordFieldMapping.IncludeTermVectors = false

	productMapping := bleve.NewDocumentMapping()

	// productMapping.AddFieldMappingsAt("sku", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("langCode", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("catalogId", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("class", englishTextFieldMapping)
	// productMapping.AddFieldMappingsAt("category", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("warrantycode", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("status", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("displayName", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("errMsg1", englishTextFieldMapping)
	// productMapping.AddFieldMappingsAt("errMsg2", keywordFieldMapping)
	// productMapping.AddFieldMappingsAt("hasImage", keywordFieldMapping)
	productMapping.AddFieldMappingsAt("discontinuedDt", dtFieldMapping)
	productMapping.AddFieldMappingsAt("launchedDt", dtFieldMapping)
	productMapping.AddFieldMappingsAt("cmUpdatedDt", dtFieldMapping)
	productMapping.AddFieldMappingsAt("cdUpdatedDt", dtFieldMapping)
	productMapping.AddFieldMappingsAt("occUpdatedDt", dtFieldMapping)
	productMapping.AddFieldMappingsAt("contentModifiedDt", dtFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("aggproduct", productMapping)

	// indexMapping.TypeField = "type"
	//indexMapping.DefaultAnalyzer = "en"
	return indexMapping, nil

}
