package ezsearch

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"

	"go.uber.org/zap/zapcore"

	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
)

type TimeSeries int

const (
	MINS30  TimeSeries = 30
	HOUR1              = 60
	HOURS3             = 3 * HOUR1
	HOURS6             = 6 * HOUR1
	HOURS12            = 12 * HOUR1
	DAY1               = 24 * HOUR1
	DAY3               = 3 * DAY1
	DAY7               = 7 * DAY1
	DAY30              = 30 * DAY1
	DAY360             = 360 * DAY1
	ALL                = 0
)

func getIndexNameWithBasepath(indexName string) string {
	if strings.Contains(indexName, "/") {
		return indexName
	}
	s := fmt.Sprintf("%s%c%s", Conf.IndexBasePath, '/', indexName)
	//fmt.Println("index path", s)
	return s
}
func getIndexNameWithTable(indexName string) string {
	if strings.Contains(indexName, "/") {
		return indexName
	}
	s := fmt.Sprintf("%s%c%s", Conf.IndexTablesPath, '/', indexName)
	//fmt.Println("index path", s)
	return s
}

//index name start with tables then that should be created under tables
// the table name pattern should be  tables.eventqueue.bleve
func GetIndex(indexName string) (bleve.Index, rest_errors.RestErr) {
	var index bleve.Index
	var err error
	//fmt.Println("Get Index called", indexName, "Count ", cache_utils.Cache.Count(), cache_utils.Cache.GetKeys())
	logger.Info(fmt.Sprintf("Get Index called %s count %d keys %v", indexName, cache_utils.Cache.Count(), cache_utils.Cache.GetKeys()))
	indexName = getIndexNameWithBasepath(indexName)
	i, err := cache_utils.Cache.Get(indexName)
	//fmt.Println("after cache get")
	if err != nil {
		logger.Error("GetIndex|Failed while cache_utils.cache", err, zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
		i, err = bleve.Open(path.Join(global.WorkingDir, indexName))
		// if i == nil || err != nil { //not available on data folder so look up in tables folder
		// 	i, err = bleve.Open(getIndexNameWithTable(indexName))
		// }
		if err != nil {
			logger.Error("GetIndex|Failed while open index", err, zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
			return nil, rest_errors.NewBadRequestError(err.Error())
		}
		//fmt.Println("GetIndex|index has found under data folder", indexName, err)
		cache_utils.AddOrUpdateCache(indexName, i)

	}

	index = i.(bleve.Index)
	//fmt.Println("before return GetIndex func ", index.Name(), cache_utils.Cache.GetKeys())
	return index, nil
}

func generateCSV(reqM *SearchRequestModel, rm *SearchResponseModel, fields map[string]int) {
	//avoid build error to import believe backage  [BuildIndexSchema|Failed schemabuild mapping no analyzer with name or type 'keyword_marker' registered]
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName
	wd, err := os.Getwd()
	if err != nil {
		logger.Error("Failed while get the working dir", err)
		return
	}
	dt := time.Now().UTC()
	p := fmt.Sprintf("%s/csvs/%s/%d%d%d", wd, reqM.IndexName, dt.Year(), dt.Month(), dt.Day())

	if _, err := os.Stat(p); os.IsNotExist(err) {
		//fmt.Println("generateCSV", err)
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			logger.Error("Failed while create file", err)
			return
		}
	}
	fn := fmt.Sprintf("products_%d_%d_%d%d%d%d.csv", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
	//fullFileName:=fmt.Sprintf("%s/%s", path, fn)
	file, err := os.Create(path.Join(p, fn))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed creating file %s", path.Join(p, fn)), err)
	}
	defer file.Close()
	if err != nil {
		logger.Error("Failed while create csv files", err)
	}
	// 2. Initialize the writer
	writer := csv.NewWriter(file)

	csvData := make([][]string, 0)

	l := len(fields)
	row := make([]string, l)
	for k, v := range fields {
		//	fmt.Println("field order", l, k, v)
		row[v] = k
	}

	//fmt.Println("print header", row)
	csvData = append(csvData, row)
	//product data
	for _, r := range rm.ResultRow {
		row = make([]string, len(fields))
		for f, val := range r {
			idx, isFound := fields[f]
			if isFound {
				//fmt.Println("col value", f, val, idx)
				row[idx] = fmt.Sprintf("%v", val)
			}
		}
		//fmt.Println("row len", len(row))
		csvData = append(csvData, row)
	}
	// 3. Write all the records
	err = writer.WriteAll(csvData) // returns error
	if err != nil {
		logger.Error("An error encountered ::", err)
	}
}

func RegisterIndexes(indexDir *string) {

	// walk the data dir and register index names
	logger.Debug("RegisterIndexes|Index path", zapcore.Field{String: *indexDir, Key: "p1", Type: zapcore.StringType})

	dirEntries, err := ioutil.ReadDir(*indexDir)
	if err != nil {
		logger.Error("RegisterIndexes|error reading data dir: %v", err)
	}

	for _, dirInfo := range dirEntries {
		indexPath := fmt.Sprintf("%v%c%s", *indexDir, os.PathSeparator, dirInfo.Name())

		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		logger.Debug("index dir info", zapcore.Field{String: dirInfo.Name(), Key: "p1", Type: zapcore.StringType})
		if !dirInfo.IsDir() {
			logger.Debug("not registering %s, skipping", zapcore.Field{String: indexPath, Key: "p1", Type: zapcore.StringType})
			continue
		}
		i, err := bleve.Open(path.Join(global.WorkingDir, indexPath))
		if err != nil {
			logger.Error("error opening index ", err, zapcore.Field{String: indexPath, Key: "p1", Type: zapcore.StringType})
		} else {
			logger.Debug("registered index:", zapcore.Field{String: dirInfo.Name(), Key: "p1", Type: zapcore.StringType})
			i.SetName(dirInfo.Name())
			cache_utils.AddOrUpdateCache(dirInfo.Name(), i)
			// set correct name in stats

		}
	}
}
