package ezsearch

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"go.etcd.io/bbolt"
	"go.uber.org/zap/zapcore"
)

func SearchResult(m SearchRequestModel) ([]*bleve.SearchResult, rest_errors.RestErr) {
	//fmt.Println("reached SearchResult|macsearchquery")
	q := bleve.NewBooleanQuery()

	if m.DateRange > 0 && len(m.DateField) > 0 {
		//fmt.Println("date range", m.DateField, m.DateRange)
		isInclusive := true
		end := time.Now().UTC()
		start := end
		start = end.Add(-time.Minute * time.Duration(m.DateRange))
		q1 := bleve.NewDateRangeInclusiveQuery(start, end, &isInclusive, &isInclusive)
		q1.SetField(m.DateField)
		q.AddMust(q1)
	}
	//fmt.Println("beofre for terms", m.Terms)
	for _, v := range m.TermQueries {

		if v == "*" {
			//fmt.Println("term query *", v)
			mq := bleve.NewMatchAllQuery()
			q.AddMust(mq)
			break
		}
		//fmt.Println("term query", v)
		vsplit := strings.Split(v, ":")
		fld := ""
		if len(vsplit) > 1 {
			fld = vsplit[0]
			v = vsplit[1]
		}
		if v == "T" {
			mq := bleve.NewBoolFieldQuery(true)
			q.AddMust(mq)

		} else if v == "F" {
			mq := bleve.NewBoolFieldQuery(false)
			q.AddMust(mq)

		} else if v[0] == '/' {
			mq := bleve.NewRegexpQuery(strings.ReplaceAll(v, "/", ""))
			if len(fld) > 0 {
				mq.SetField(fld)
			}
			q.AddMust(mq)

		} else {
			mq := bleve.NewMatchQuery(v)
			if len(fld) > 0 {
				mq.SetField(fld)
			}
			q.AddMust(mq)
		}
	}
	if len(m.PharseQueries) > 0 {
		for _, v := range m.PharseQueries {
			if v == "*" {
				mq := bleve.NewMatchAllQuery()
				q.AddMust(mq)
				break
			}
			vsplit := strings.Split(v, ":")
			fld := ""
			if len(vsplit) > 1 {
				fld = vsplit[0]
				v = vsplit[1]
			}
			if v[0] == '/' {
				mq := bleve.NewRegexpQuery(strings.ReplaceAll(v, "/", ""))
				if len(fld) > 0 {
					mq.SetField(fld)
				}
				q.AddMust(mq)

			} else {

				mpq := bleve.NewMatchPhraseQuery(v)
				if len(fld) > 0 {
					mpq.SetField(fld)
				}
				q.AddMust(mpq)
			}
		}
	}

	req := bleve.NewSearchRequest(q)

	if len(m.Fields) > 0 {
		req.Fields = m.Fields
	} else {
		req.Fields = []string{"*"}
	}
	if len(m.SortBy) > 0 {
		req.SortBy(m.SortBy)
	}
	for _, f := range m.Facets {
		req.AddFacet(f, bleve.NewFacetRequest(f, 10))
	}
	req.From = m.From
	req.Size = m.Size
	if m.Size == 0 {
		req.Size = 50
	}
	//fmt.Println("index name", m.IndexName)
	//fmt.Println("index names", names)
	//fmt.Println("index loop SearchResult|macsearchquery", indexName)
	//	fmt.Println("index name", indexName, index.Name())
	//return nil, rest_errors.NewInternalServerError("Failed while execute the query", err)
	results := exeSearch(strings.Split(m.IndexName, "|"), m.DateRange, req)
	//fmt.Println("end of func SearchResult|macsearchquery")
	return results, nil

}

func exeSearch(indexNames []string, duration int32, req *bleve.SearchRequest) []*bleve.SearchResult {
	//names := getIndexNameUsingDt(indexNames, duration)

	results := make([]*bleve.SearchResult, 0)

	for _, indexName := range indexNames {
		//fmt.Println("exeSearch|indexname", indexName)
		index, err := GetIndex(indexName)
		if err != nil {
			logger.Error("Failed while get index", err, zapcore.Field{String: indexName, Type: zapcore.StringType})
			continue
		}
		//logger.Debug("exeSearch|index names", zapcore.Field{String: index.Name(), Key: "p1", Type: zapcore.StringType})

		res, err1 := index.Search(req)
		if err1 != nil {
			logger.Error("Failed while execute the query", err1)
			continue
		}
		if res.Hits != nil && len(res.Hits) > 0 {
			//	fmt.Println("exeSearch|index names", indexName, res.Hits.Len())
			// if res.Facets != nil && len(res.Facets) > 0 {
			// 	fmt.Println("merge facets", len(res.Facets), res.Facets)
			// 	for _, f := range res.Facets {
			// 		for _, v := range f.Terms {
			// 			fmt.Println("merge facets Term", v.Term, v.Count)
			// 		}
			// 	}
			// }

			results = append(results, res)
		}

	}
	return results
}
func CheckExpirationCallback(key string, value interface{}) bool {
	//db, err := bolt.Open("my.db", 0600, nil)
	switch value.(type) {
	case bleve.Index:
		value.(bleve.Index).Close()
	case *bbolt.DB:
		value.(*bbolt.DB).Close()
	default:
		logger.Info("CheckExpirationCallback type not found!", zapcore.Field{String: key, Key: "p1", Type: zapcore.StringType})
	}
	return true
}

func getIndexNameUsingDt(indexNames []string, startDt time.Time, endDt time.Time) []string {
	//fmt.Println("getIndexNameUsingDt|start", indexNames)
	names := make([]string, 0)
	for _, indexName := range indexNames {
		grp := global.RegexParseDate.FindAllSubmatch([]byte(indexName), -1)
		if grp == nil {
			names = append(names, indexName)
			continue
		}
		dtFormat := string(grp[0][1])
		//fmt.Println("getIndexNameUsingDt", dtFormat, startDt, endDt)
		switch len(dtFormat) {
		case 10: //day wise index name
			names = append(names, getRangeIndexes(startDt, dtFormat, indexName, endDt, 0, 0, 1)...)
		case 7:
			names = append(names, getRangeIndexes(startDt, dtFormat, indexName, endDt, 0, 1, 0)...)
		default:
			//needs to be revisited sime date := time.Date(2013,time.January,31,23,59,59,0,time.UTC) jumps two months
			names = append(names, getRangeIndexes(startDt, dtFormat, indexName, endDt, 1, 0, 0)...)
		}
	}
	//fmt.Println("getIndexNameUsingDt|end", names)
	return names

}

func getRangeIndexes(startDt time.Time, dtFormat string, indexName string, endDt time.Time, y int, m int, d int) []string {
	curDt := startDt
	//fmt.Println("getRangeIndexes", curDt.Format(dtFormat), endDt.Format(dtFormat))
	names := make([]string, 0)
	for {
		dtVal := curDt.Format(dtFormat)
		name := strings.Replace(indexName, fmt.Sprintf("{%s}", dtFormat), dtVal, -1)
		names = append(names, name)
		//fmt.Println("getRangeIndexes", curDt.Format(dtFormat), endDt.Format(dtFormat))
		if curDt.Format(dtFormat) == endDt.Format(dtFormat) {
			break
		}
		curDt = curDt.AddDate(y, m, d)
	}
	return names
}

func MergeSearchResult(searchResults []*bleve.SearchResult) (*SearchResponseModel, rest_errors.RestErr) {

	if len(searchResults) == 0 {
		return nil, nil
	}
	res := searchResults[0]

	for idx, r := range searchResults {
		if idx == 0 {
			continue
		}
		res.Merge(r)
	}

	rm := SearchResponseModel{}
	if res.Hits.Len() > 0 {
		rm.ResultRow = make([]map[string]interface{}, 0)
	}
	rm.Fields = res.Request.Fields
	// resS, _ := json.Marshal(res)
	// fmt.Println("search res", string(resS))
	rm.Total = res.Total
	rm.Took = res.Took
	rm.Status = *res.Status
	for _, rv := range res.Hits {
		rm.ResultRow = append(rm.ResultRow, rv.Fields)
	}
	if len(res.Facets) > 0 {
		rm.Facets = make(map[string][]EzTermFacet)
	}
	for k, v := range res.Facets {
		terms := make([]EzTermFacet, 0)
		for _, fv := range v.Terms {
			terms = append(terms, EzTermFacet{Term: fv.Term, Count: fv.Count})
		}
		rm.Facets[k] = terms
	}
	// resbyte, err := json.Marshal(rm)
	// if err != nil {
	// 	fmt.Println("Failed while marshal final search result", err)
	// 	return nil, rest_errors.NewRestError("Failed while final search result", http.StatusInternalServerError, err.Error(), nil)
	// }
	//resStr := string(resbyte)
	//fmt.Println("Final search", resStr)
	//header
	if len(rm.ResultRow) == 0 {
		return nil, nil
	}
	fieldPos := make(map[string]int)
	if rm.Fields[0] != "*" {
		logger.Debug("read from fields list as noticed *")
		for idx, fn := range rm.Fields {
			fieldPos[fn] = idx
		}
	} else {
		idx := 0
		for f := range rm.ResultRow[0] {
			fieldPos[f] = idx
			idx = idx + 1
		}
	}

	//generateCSV(&m, &rm, fieldPos)
	return &rm, nil
}
