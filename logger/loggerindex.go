package logger

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/token/keyword"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
)

// the table name pattern should be  tables.eventqueue.bleve
func getIndex() (bleve.Index, rest_errors.RestErr) {
	var (
		index bleve.Index
		err   error
	)
	lIndexName := Conf.ApplogIndexPath

	if hasIndexDatePattern {
		dtStr := time.Now().UTC().Format(loggerPatternName)
		lIndexName = strings.Replace(Conf.ApplogIndexPath, fmt.Sprintf("{%s}", loggerPatternName), dtStr, -1)

	}
	//fmt.Println("lIndexName", lIndexName)
	//fmt.Println("Get Index called", ApplogIndexPath, "Count ", cache_utils.Cache.GetKeys())
	i, err := cache_utils.Cache.Get(lIndexName)
	//fmt.Println("after cache get")
	if err != nil || i == nil {
		//fmt.Println("logger|getIndex|Failed while cache_utils.cache", ApplogIndexPath, err)
		i, err = bleve.Open(path.Join(global.WorkingDir, lIndexName))
		//fmt.Println("#29|logger|getIndex|Failed while cache_utils.cache", lIndexName, err)
		if err != nil {
			//fmt.Println("logger|getIndex|logger|Failed while open index", lIndexName, err)
			BuildAppIndexSchema() //after created index from build schema would closed automatically so need to open
			i, _ = bleve.Open(path.Join(global.WorkingDir, lIndexName))
		}
		//fmt.Println("GetIndex|index has found under data folder", indexName, err)
		cache_utils.AddOrUpdateCache(lIndexName, i)
		//fmt.Println("#34|logger|getIndex|Failed while cache_utils.cache", cache_utils.Cache.GetKeys())
	}
	index = i.(bleve.Index)
	//fmt.Println("logger|getIndex|", index.Name())
	return index, nil
}

type BleveFieldDef struct {
	Name string `json:"name"`
	Type string `json:"type"` //possible values [bool|text|date|numeric|geopoint]
}

func BuildAppIndexSchema() rest_errors.RestErr {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	indexmapping := bleve.NewDocumentMapping()
	indexmapping.AddFieldMappingsAt("l", englishTextFieldMapping)
	indexmapping.AddFieldMappingsAt("t", englishTextFieldMapping)
	indexmapping.AddFieldMappingsAt("m", englishTextFieldMapping)
	fields := []BleveFieldDef{{Name: "l", Type: "text"}, {Name: "t", Type: "date"}, {Name: "m", Type: "text"}}

	indexMapping := bleve.NewIndexMapping()
	docMapName := "docs" //strings.ReplaceAll(indexName, ".", "")
	indexMapping.AddDocumentMapping(docMapName, indexmapping)
	var lIndexName string
	if hasIndexDatePattern {
		dtStr := time.Now().UTC().Format(loggerPatternName)
		lIndexName = strings.Replace(Conf.ApplogIndexPath, fmt.Sprintf("{%s}", loggerPatternName), dtStr, -1)

	}
	index, err := bleve.Open(path.Join(global.WorkingDir, lIndexName))
	//fmt.Println("afterGetIndex(indexName)")
	if err != nil || index == nil {

		//fmt.Println("logger|BuildIndexSchema|Creating  new index ... ", lIndexName)
		// create a mapping

		index, err := bleve.New(lIndexName, indexMapping)
		if err != nil {
			fmt.Println("BuildIndexSchema|Failed product index field mapping", err)
			return rest_errors.NewBadRequestError(err.Error())
		}
		//cache_utils.AddOrUpdateCache(ApplogIndexPath, index)
		index.Close()
	} else {
		key := fmt.Sprintf("%s.schema", Conf.ApplogIndexPath)
		bytes, _ := coredb.GetKey(key)
		//fmt.Println("before if not found key in core db", key, string(bytes))
		if bytes == nil || len(bytes) == 0 {
			fmt.Println("not found key in core db", key, string(bytes))
			bytes, _ := json.Marshal(fields)
			coredb.AddKey(key, bytes)
		}
		if index != nil {
			index.Close()
		}

		//return rest_errors.NewInternalServerError(errors.New(msg))
	}
	bytes, _ := json.Marshal(fields)
	coredb.AddKey(fmt.Sprintf("%s.schema", Conf.ApplogIndexPath), bytes)
	return nil
}
