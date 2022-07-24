package ezsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"go.uber.org/zap/zapcore"
)

var (
	since map[string]int
)

func init() {
	since = make(map[string]int)
	since["seconds"] = 1
	since["minutes"] = 1 * 60
	since["minute"] = 1 * 60
	since["hour"] = 1 * 60 * 60
	since["hours"] = 1 * 60 * 60
	since["day"] = 1 * 60 * 60 * 24
	since["days"] = 1 * 60 * 60 * 24
}

func GetFields(indexName string) []string {
	i, err := GetIndex(indexName)
	if i == nil || err != nil {
		logger.Warn("Index is missing ", zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
		return nil
	}

	fields, _ := i.Fields()
	//fmt.Println("controllers|getfields", indexName, fields)
	return fields
}
func GetValues(indexName string, fieldName string) []interface{} {
	// //i := IndexNameMapping[indexName]
	// if i == nil {
	// 	println(fmt.Sprintf("Index is missing %s", indexName))
	// 	return nil
	// }
	//srm := SearchRequestModel{Fields: []string{fieldName}, IndexName: indexName, From: 0, Size: 1}
	queryStr := fmt.Sprintf("select %s from %s", fieldName, indexName)
	result, err := PostSearchResult(queryStr)
	if err != nil {
		logger.Error("Failed while fetch Getvalues", err, zapcore.Field{String: indexName, Type: zapcore.StringType}, zapcore.Field{String: fieldName, Key: "p1", Type: zapcore.StringType})
		return nil
	}
	list := make([]interface{}, 0)
	for _, row := range result.ResultRow {
		list = append(list, row[fieldName])
	}
	return list
}

// func AddOrUpdateIndex(indexName string, docId string, model interface{}) error {
// 	i := GetIndex(indexName)

// 	if i == nil {
// 		println(fmt.Sprintf("Index is missing %s", indexName))
// 		CreateNewIndex(indexName, []string{"timestamp"})
// 		//return fmt.Errorf("index is missing %s", indexName)
// 	}
// 	fmt.Println("Index name", docId, indexName, model)
// 	err := i.Index(docId, model)
// 	return err
// }

// func AddBatchIndex(indexName string, docIds []string, models []interface{}) error {
// 	i := GetIndex(indexName)
// 	batch := i.NewBatch()
// 	for idx, id := range docIds {
// 		batch.Index(id, models[idx])
// 	}
// 	return i.Batch(batch)
// }

// func DeleteIndex(indexName string, docId string, model interface{}) error {
// 	i := GetIndex(indexName)
// 	if i == nil {
// 		println(fmt.Sprintf("Index is missing %s", indexName))
// 		return fmt.Errorf("index is missing %s", indexName)
// 	}
// 	return i.Delete(docId)
// }

func PostSearchResult(queryStr string) (*SearchResponseModel, rest_errors.RestErr) {

	req, indexNames, err := parseQuery(queryStr)
	if err != nil {
		errStr := fmt.Sprintf("Failed while parse query %s %s", indexNames, err.Error())
		logger.Error(errStr, err)
		return nil, rest_errors.NewBadRequestError(errStr)
	}

	str, _ := json.Marshal(req)
	//fmt.Println("qery")

	logger.Debug("PostSearchResult|req model", zapcore.Field{String: string(str), Key: "p1", Type: zapcore.StringType})

	resColl := exeSearch(indexNames, 0, req)
	if resColl == nil || len(resColl) == 0 {
		msg := fmt.Sprintf("No response from execute search %s", queryStr)
		logger.Error(msg, errors.New(""))
		return nil, rest_errors.NewNotFoundError(msg)
	}
	res := resColl[0]
	// fmt.Println("rescoll count", len(resColl))
	// fmt.Println("res res.Hits.Len() ", res.Hits.Len())

	if len(resColl) > 1 {
		for i := 1; i < len(resColl); i++ {
			//fmt.Println("merge response", resColl[i])
			res.Merge(resColl[i])

		}
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
	//println("response collection", len(resColl), res.Hits)
	for _, rv := range res.Hits {

		rm.ResultRow = append(rm.ResultRow, rv.Fields)
		//str, _ := json.Marshal(rv.Fields)
		//println("row data", len(rv.Fields), string(str))
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
	//resbyte, err := json.Marshal(rm)
	if err != nil {
		logger.Error("Failed while marshal final search result", err)
		return nil, rest_errors.NewRestError("Failed while final search result", http.StatusInternalServerError, err.Error(), nil)
	}
	//resStr := string(resbyte)
	//fmt.Println("Final search", resStr)
	//header
	if len(rm.ResultRow) == 0 {
		return nil, nil
	}
	// fieldPos := make(map[string]int)
	// if len(rm.Fields) > 0 && rm.Fields[0] != "*" {
	// 	//fmt.Println("read from fields list")
	// 	for idx, fn := range rm.Fields {
	// 		fieldPos[fn] = idx
	// 	}
	// } else {
	// 	idx := 0
	// 	for f := range rm.ResultRow[0] {
	// 		fieldPos[f] = idx
	// 		idx = idx + 1
	// 	}
	// }

	//generateCSV(&m, &rm, fieldPos)
	return &rm, nil
}

func constructQueries(qList *[]*whereClause, q *query.BooleanQuery) string {
	logger.Debug(fmt.Sprintf("constructQueries "), zapcore.Field{Integer: int64(len(*qList)), Type: zapcore.Int64Type})

	if len(*qList) == 0 {
		logger.Debug("constructQueries|qList == nil all queris")
		q.AddShould(bleve.NewMatchAllQuery())
		return ""
	}
	for _, fd := range *qList {
		var genQry query.Query
		var err error
		if strings.Contains(fd.Min, "*") || strings.Contains(fd.Min, "?") {
			logger.Debug("wild card")
			gq := query.NewWildcardQuery(fd.Min)
			if len(fd.Field) > 0 {
				gq.SetField(fd.Field)
			}
			genQry = gq
		} else if len(fd.Min) > 3 && fd.Min[0] == '/' { //consider regex
			logger.Debug("regex query")
			gq := query.NewRegexpQuery(strings.ReplaceAll(fd.Min, "/", ""))
			if len(fd.Field) > 0 {
				gq.SetField(fd.Field)
			}
			genQry = gq
		} else if len(fd.Field) == 0 { //consider pharse
			gq := query.NewPrefixQuery(fd.Min)
			fd.MustMustNotShould = 1
			genQry = gq
		} else {
			switch fd.BleveQueryType {
			case BooleanQuery:
				logger.Debug("constructQueries|bool", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType})
				b := true
				if strings.ToLower(fd.Min) == "false" {
					b = false
				}
				gq := query.NewBoolFieldQuery(b)
				gq.SetField(fd.Field)
				//fmt.Println(gq.FieldVal, fd.Field, gq.Bool, b)
			case MatchQuery:
				logger.Debug("constructQueries|Matchquery", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType})
				gq := query.NewMatchQuery(fd.Min)
				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq
			case MatchPhraseQuery:
				logger.Debug("constructQueries|matchPharse", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType})
				gq := query.NewMatchPhraseQuery(fd.Min)
				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq
			case DateRangeInclusiveQuery:
				logger.Debug("constructQueries|DateRangeInclusiveQuery", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType}, zapcore.Field{String: fd.Max, Key: "max", Type: zapcore.StringType})
				var mDt, dt time.Time
				if len(fd.Min) > 0 {
					mDt, err = time.Parse(date_utils.UTCDateLayout, fd.Min)
					if err != nil {
						return err.Error()
					}
				}
				if len(fd.Max) > 0 {
					dt, err = time.Parse(date_utils.UTCDateLayout, fd.Max)
					if err != nil {
						return err.Error()
					}
				}
				gq := query.NewDateRangeInclusiveQuery(dt, mDt, &fd.IsMaxInclusive, &fd.IsMinInclusive)

				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq
			case DateRangeQuery:
				logger.Debug(fmt.Sprintf("constructQueries|DateRangeQuery min %s, max %s", fd.Min, fd.Max))
				var mDt, dt time.Time
				if len(fd.Min) > 0 {
					mDt, err = time.Parse(date_utils.UTCDateLayout, fd.Min)
					if err != nil {
						return err.Error()
					}
				}
				if len(fd.Max) > 0 {
					dt, err = time.Parse(date_utils.UTCDateLayout, fd.Max)
					if err != nil {
						return err.Error()
					}

				}
				gq := query.NewDateRangeQuery(dt, mDt)
				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq
			case NumericRangeInclusiveQuery:
				logger.Debug("constructQueries|NumericRangeInclusiveQuery", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType}, zapcore.Field{String: fd.Max, Key: "max", Type: zapcore.StringType})
				var mi *float64
				var mx *float64
				if len(fd.Min) > 0 {
					m, err := strconv.ParseFloat(fd.Min, 32)

					if err != nil {
						return err.Error()
					}
					mi = &m
				}
				if len(fd.Max) > 0 {
					m, err := strconv.ParseFloat(fd.Max, 32)
					if err != nil {
						return err.Error()
					}
					mx = &m
				}
				gq := query.NewNumericRangeInclusiveQuery(mx, mi, &fd.IsMaxInclusive, &fd.IsMinInclusive)
				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq
			case TermRangeQuery:
				logger.Debug("constructQueries|TermRangeQuery", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType}, zapcore.Field{String: fd.Max, Key: "max", Type: zapcore.StringType})

				gq := query.NewTermRangeInclusiveQuery(fd.Min, fd.Max, &fd.IsMaxInclusive, &fd.IsMinInclusive)
				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq
			case NumericRangeQuery:
				logger.Debug("constructQueries|NumericRangeQuery", zapcore.Field{String: fd.Min, Key: "min", Type: zapcore.StringType}, zapcore.Field{String: fd.Max, Key: "max", Type: zapcore.StringType})

				var mi *float64
				var mx *float64
				if len(fd.Min) > 0 {
					m, err := strconv.ParseFloat(fd.Min, 32)
					if err != nil {
						return err.Error()
					}
					mi = &m
				}

				if len(fd.Max) > 0 {
					m, err := strconv.ParseFloat(fd.Max, 32)
					if err != nil {
						return err.Error()
					}
					mx = &m
				}
				gq := query.NewNumericRangeQuery(mx, mi)
				if len(fd.Field) > 0 {
					gq.SetField(fd.Field)
				}
				genQry = gq

			}
		}
		//fmt.Println("fd.ReqOptExc", genQry)
		if fd.MustMustNotShould == 1 {
			q.AddMust(genQry)
		} else if fd.MustMustNotShould == 2 {
			q.AddMustNot(genQry)
		} else {
			q.AddShould(genQry)
		}
	}
	return ""
}

type queryDef struct {
	indexName      string
	from           int
	size           int
	whereCondtions []whereClause
}
type whereClause struct {
	Field             string
	Min               string
	Max               string
	IsMinInclusive    bool
	IsMaxInclusive    bool
	hasRangeOperator  bool
	MustMustNotShould int8 //"+" MUST-->1 "-" MUST NOT-->2 without these SHOULD-->0 clause
	BleveQueryType    queryType
}
type queryType int8

const (
	MatchPhraseQuery queryType = iota
	BooleanQuery
	DateRangeQuery
	DateRangeInclusiveQuery
	NumericRangeQuery
	NumericRangeInclusiveQuery
	MatchQuery
	GeoPointQuery
	TermRangeQuery
)

func parseSinceClause(sinceStr string) (string, map[string][]string) {
	split := strings.Split(sinceStr, "since")
	if len(split) == 1 {
		return sinceStr, nil
	}
	clauses := make(map[string][]string)

	sinceSplit := (strings.Split(strings.Trim(split[1], " "), " "))
	words := make([]string, 0)
	for _, v := range sinceSplit {
		tv := strings.Trim(v, " ")

		if len(tv) > 0 {
			words = append(words, tv)
		}
	}
	clauses["sin"] = sinceSplit
	return split[0], clauses

}
func parseQuery(parseVal string) (*bleve.SearchRequest, []string, error) {
	//logger.Info(parseVal)

	hasSinceClause := true
	hasDatePattern := false
	parseVal = strings.ReplaceAll(parseVal, `\n`, " ") //remove new line
	parseVal = strings.ReplaceAll(parseVal, `  `, " ")
	parseVal = strings.Trim(parseVal, " ")
	//fmt.Println("testing parseVal", parseVal)
	kw := getParsedQueryByKeyword(parseVal)
	//parseStr, clauses := parseSinceClause(parseVal)

	//fmt.Println(parseStr, parseVal, clauses)

	hasRangeOperator := false
	if strings.Contains(parseVal, ">") {
		hasRangeOperator = true
	}
	if !hasRangeOperator {
		if strings.Contains(parseVal, "<") {
			hasRangeOperator = true
		}
	}
	sel := kw["sel"]
	if len(sel) == 0 {
		errStr := fmt.Sprintf("query does not have select and should be small letter %s", parseVal)
		logger.Error(errStr, errors.New(""))
		return nil, nil, errors.New(errStr)
	}
	//s := `select 2,3,3,4 from indexName|index2 where +name:test, +startDt:>"2016-2-2",startDt:<="2016-2-2" sort -name faces:cateory:20,brand:5 limit 0,70 facets name`

	fro := kw["fro"]
	if len(fro) == 0 {
		errStr := fmt.Sprintf(`%s unable to get index name in from clause or from clause is missing.`, parseVal)
		logger.Error(errStr, errors.New(""))
		return nil, nil, errors.New(errStr)
	}
	indexNames := strings.Split(fro, "|")
	var schemaDefs []common.BleveFieldDef
	schemaDefs, _ = GetBleveIndexSchema(indexNames[0])

	indexes := common.GetAllIndexes()
	//fmt.Println(indexes)
	found := false
	indexList := make([]string, 0)
	for _, indexName := range indexNames {

		hasDatePattern = strings.Contains(indexName, "{")
		if !strings.Contains(indexName, "/") {
			indexName = getIndexNameWithBasepath(indexName)
		}
		_, found = indexes[indexName]
		if found || hasDatePattern { //date pattern is provided then skip index exist logic
			indexList = append(indexList, indexName)
			found = true
			//break
		}
	}

	var fdList []*whereClause
	//since clause
	var sinceWhereCond *whereClause
	var err error
	sin := kw["sin"]
	if hasDatePattern && len(sin) == 0 {
		errStr := fmt.Sprintf(`query has date pattern but missing since caluse [%s]`, parseVal)
		logger.Error(errStr, errors.New(""))
		return nil, nil, errors.New(errStr)
	}

	if len(sin) > 0 {
		hasSinceClause = true
		sinSplit := strings.Split(sin, " ")
		sinList := make([]string, 0)
		for _, v := range sinSplit {
			if len(v) > 0 {
				sinList = append(sinList, strings.Trim(v, " "))
			}
		}
		sinceWhereCond, err = parseSinceField(sinList)
		if err != nil {
			return nil, nil, err
		}

		fdList = append(fdList, sinceWhereCond)
	}

	if !found {
		errStr := fmt.Sprintf(`%v is/are not available`, indexList)
		logger.Error(errStr, errors.New(""))
		return nil, nil, errors.New(errStr)
	}
	indexNames = indexList
	q := bleve.NewBooleanQuery()
	req := bleve.NewSearchRequest(q)
	for _, field := range strings.Split(sel, ",") {
		if field == "*" {
			if len(schemaDefs) == 0 {
				req.Fields = GetFields(indexNames[0])
			} else {
				for _, v := range schemaDefs {
					req.Fields = append(req.Fields, v.Name)
				}
			}
		} else {
			req.Fields = append(req.Fields, strings.Split(strings.Trim(field, " "), ",")...)
		}
	}
	//set wild card if there is no schema found from core db
	if len(req.Fields) == 0 {
		req.Fields = []string{"*"}
	}

	fac := kw["fac"]
	for _, field := range strings.Split(fac, ",") {
		split := strings.Split(field, ",")
		for _, f := range split {
			fs := strings.Split(strings.Trim(f, " "), ":")
			c := 20
			if len(fs) > 1 {
				c, _ = strconv.Atoi(fs[1])
			}
			req.AddFacet(fs[0], bleve.NewFacetRequest(fs[0], c))
		}
	}
	lim := kw["lim"]

	if len(lim) == 0 {
		req.From = 0
		req.Size = 50
	} else {

		fs := strings.Split(lim, ",")
		logger.Debug("limit", zapcore.Field{String: fmt.Sprintf("%v", fs), Key: "p1", Type: zapcore.StringType})
		if len(fs) == 1 {
			req.From, _ = strconv.Atoi(fs[0])
		} else {
			req.From, _ = strconv.Atoi(fs[0])
			req.Size, _ = strconv.Atoi(fs[1])
		}
	}

	if req.Size == 0 {
		req.Size = 50
	}
	sortList := make([]string, 0)
	sor := kw["sor"]

	sortList = append(sortList, strings.Split(sor, ",")...)
	req.SortBy(sortList)

	whe := kw["whe"]
	split := strings.Split(whe, ",")
	for _, v := range split {
		//fmt.Println("where for loop", v)
		if len(v) > 0 {
			parseQueryField(strings.Trim(v, " "), &fdList, schemaDefs)
		}
	}
	errmsg := constructQueries(&fdList, q)
	if len(errmsg) > 0 {
		return nil, nil, rest_errors.NewBadRequestError(errmsg)
	}

	if !hasSinceClause && hasDatePattern {
		logger.Error(fmt.Sprintf("Pattern index %v mush have since clause", indexNames), errors.New(""))
		return nil, nil, rest_errors.NewBadRequestError(fmt.Sprintf("Pattern index %v mush have since clause", indexNames))
	}
	if sinceWhereCond != nil { //get the list of indexes when use pattern along with since clause must be available
		stDt, _ := time.Parse(date_utils.UTCDateLayout, sinceWhereCond.Max)
		endDt, _ := time.Parse(date_utils.UTCDateLayout, sinceWhereCond.Min)
		list := getIndexNameUsingDt(indexNames, stDt, endDt)
		indexNames = make([]string, 0)
		for _, n := range list {
			if indexes[n] {
				indexNames = append(indexNames, n)
			}

		}
	}

	return req, indexNames, nil
}

func findWhereField(fieldName string, list *[]*whereClause) (*whereClause, bool) {
	var item *whereClause
	fieldName = strings.ReplaceAll(fieldName, "+", "")
	fieldName = strings.ReplaceAll(fieldName, "-", "")
	for _, item = range *list {
		// fmt.Println("findWhereField|for", fieldName, item)
		if item.Field == fieldName {
			// fmt.Println("findWhereField|for|return", fieldName, item)
			return item, true
		}
	}
	return item, false
}

func findFieldFromSchema(fieldName string, list []common.BleveFieldDef) *common.BleveFieldDef {
	for _, item := range list {
		if item.Name == fieldName {
			return &item
		}
	}
	return nil
}

func parseSinceField(sinceClause []string) (*whereClause, error) {
	logger.Debug("sincecaluse")
	if len(sinceClause) < 3 {
		logger.Error(fmt.Sprintf("Since is not well constructed must be 3 words after since] %v", sinceClause), errors.New(""))
		return nil, errors.New("")
	}

	wc := whereClause{}
	fieldsplit := strings.Split(sinceClause[0], ":")
	if len(fieldsplit) == 1 {
		logger.Error(fmt.Sprintf("Since clause field name is missing %s", sinceClause[0]), errors.New(""))
		return nil, errors.New(fmt.Sprintf("Since clause field name is missing %s", sinceClause[0]))
	}

	wc.Field = fieldsplit[0]
	dur, err := strconv.Atoi(fieldsplit[1])
	if err != nil {
		return nil, err
	}
	eDt := time.Now().UTC()
	if strings.ToLower(sinceClause[2]) != "ago" {
		eDt, err = time.Parse(date_utils.UTCDateLayout, sinceClause[2])
		if err != nil {
			return nil, err
		}
	}
	var minutes int
	var ok bool

	if minutes, ok = since[strings.ToLower(sinceClause[1])]; !ok {
		logger.Error(fmt.Sprintf("Since clause is not set properly %s", sinceClause[1]), errors.New(""))
		return nil, errors.New(fmt.Sprintf("Since clause is not set properly %s", sinceClause[1]))
	}
	//logger.Debug("since nown", zapcore.Field{String: eDt.Format(date_utils.UTCDateLayout), Type: zapcore.StringType}, zapcore.Field{String: sinceClause[1], Type: zapcore.StringType}, zapcore.Field{Integer: int64(minutes), Type: zapcore.Int64Type})
	sDt := eDt.Add(-time.Second * time.Duration(minutes*dur))
	//logger.Debug("since nown", zapcore.Field{Integer: int64(-time.Minute * time.Duration(minutes*dur)), Type: zapcore.Int64Type}, zapcore.Field{String: sDt.Format(date_utils.UTCDateLayout), Type: zapcore.StringType}, zapcore.Field{String: eDt.Format(date_utils.UTCDateLayout), Type: zapcore.StringType}, zapcore.Field{String: sinceClause[1], Type: zapcore.StringType}, zapcore.Field{Integer: int64(minutes), Type: zapcore.Int64Type})
	wc.IsMaxInclusive = true
	wc.IsMinInclusive = true
	wc.Max = sDt.Format(date_utils.UTCDateLayout)
	wc.Min = eDt.Format(date_utils.UTCDateLayout)
	wc.BleveQueryType = DateRangeInclusiveQuery
	wc.MustMustNotShould = 1
	//fmt.Println("since field def", wc)
	return &wc, nil

}
func parseQueryField(fieldVal string, list *[]*whereClause, schemaDefs []common.BleveFieldDef) {
	var field, val string
	val = fieldVal
	for idx := 0; idx < len(fieldVal); idx++ {
		if fieldVal[idx] == ':' {
			field = fieldVal[0:idx]
			val = fieldVal[idx+1:]
			break
		}
	}

	var rf *whereClause
	var exist bool = false

	if strings.Contains(fieldVal, "<") || strings.Contains(fieldVal, ">") { //construct range fields min and max
		//fmt.Println("strings.Contains(fieldVal, <)", fieldVal, len(*list))
		rf, exist = findWhereField(field, list)

	}
	if !exist {
		//fmt.Println("new|else!esit", fieldVal)
		rf = &whereClause{}
	}
	if len(val) > 0 {
		//set field eliminte - symbol
		rf.Field, rf.MustMustNotShould = getReqOptInc(field)
	}

	setRangeField(rf, val)
	setQueryType(rf, schemaDefs)
	if !exist {
		//fmt.Println("inside if!esit", rf)
		*list = append(*list, rf)
	}
	// } else {

	// 	fmt.Println("else!esit", rf)
	// }
	//fmt.Println("parsequeryField|before return|", len(*list))
	return
}

func setQueryType(rf *whereClause, list []common.BleveFieldDef) {
	sfd := findFieldFromSchema(rf.Field, list)
	//logger.Debug(fmt.Sprintf("setQueryType|min:%s max:%s", rf.Min, rf.Max))
	if len(strings.Split(rf.Min, " ")) > 1 {
		rf.BleveQueryType = MatchPhraseQuery
	} else {
		rf.BleveQueryType = MatchQuery
	}
	if sfd != nil {
		switch sfd.Type {
		case "date":
			if rf.IsMaxInclusive || rf.IsMinInclusive {
				rf.BleveQueryType = DateRangeInclusiveQuery
			} else if rf.hasRangeOperator {
				rf.BleveQueryType = DateRangeQuery
			} else {
				rf.BleveQueryType = DateRangeInclusiveQuery
				if len(rf.Min) == 0 {
					rf.Min = rf.Max
				}
				if len(rf.Max) == 0 {
					rf.Max = rf.Min
				}
				rf.IsMinInclusive = true
				rf.IsMaxInclusive = true
				//fmt.Println("case|date|else", rf)
			}

		case "numeric":
			if rf.IsMaxInclusive || rf.IsMinInclusive {
				rf.BleveQueryType = NumericRangeInclusiveQuery
			} else if rf.hasRangeOperator {
				rf.BleveQueryType = NumericRangeQuery
				//logger.Debug(fmt.Sprintf("numberic  %s %s, %s %v", rf.Field, rf.Min, rf.Max, rf))
			} else {
				rf.BleveQueryType = NumericRangeInclusiveQuery
				if len(rf.Min) == 0 {
					rf.Min = rf.Max
				}
				if len(rf.Max) == 0 {
					rf.Max = rf.Min
				}
				rf.IsMinInclusive = true
				rf.IsMaxInclusive = true
			}

		case "bool":
			rf.BleveQueryType = BooleanQuery
		case "geo":
			rf.BleveQueryType = GeoPointQuery
		}

	} else if len(rf.Field) > 0 && len(rf.Min) > 0 && len(rf.Max) > 0 {
		rf.BleveQueryType = TermRangeQuery
	}
	//fmt.Println("set query Type", rf)
}
func getReqOptInc(field string) (string, int8) {
	if len(field) == 0 {
		return "", 0
	}
	retVal := field
	var optVal int8
	optVal = 0

	if field[0] == '+' { //eliminate + symbol
		optVal = 1
		retVal = field[1:]
	} else if field[0] == '-' {
		retVal = field[1:] //eliminate - symbol
		optVal = 2
	}
	return retVal, optVal
}
func setRangeField(rf *whereClause, fieldVal string) {
	fieldVal = strings.ReplaceAll(fieldVal, `"`, "")
	var op string
	if len(fieldVal) > 1 && fieldVal[0] >= 60 && fieldVal[0] <= 62 { //ascii value for <=> symbols
		op = string(fieldVal[0])
		fieldVal = fieldVal[1:]
	}
	//fmt.Println("fieldVal", op, len(fieldVal), fieldVal)
	if len(fieldVal) > 1 && fieldVal[0] >= 60 && fieldVal[0] <= 62 {
		op = fmt.Sprintf("%s%c", op, fieldVal[0])
		fieldVal = fieldVal[1:]
	}
	//fmt.Println("fieldVal", op, len(fieldVal), fieldVal)
	// rf.IsMaxInclusive = false
	// rf.IsMinInclusive = false
	switch op {
	case ">":
		rf.hasRangeOperator = true
		rf.Max = fieldVal
		rf.IsMaxInclusive = false
	case ">=":
		rf.hasRangeOperator = true
		rf.Max = fieldVal
		rf.IsMaxInclusive = true
	case "<":
		rf.hasRangeOperator = true
		rf.Min = fieldVal
		rf.IsMaxInclusive = false
	case "<=":
		rf.hasRangeOperator = true
		rf.Min = fieldVal
		rf.IsMinInclusive = true
	default:
		rf.Min = fieldVal
	}
	//fmt.Println(fmt.Sprintf("setRangeField|after Op %s fieldValue=%s, min=%s %t,  max=%s %t", op, fieldVal, rf.Min, rf.IsMinInclusive, rf.Max, rf.IsMaxInclusive))
}

func GetTimeStampFieldValue(model interface{}) interface{} {
	const fld = "timestamp"
	getType := reflect.TypeOf(model)
	getValue := reflect.ValueOf(model)
	for i := 0; i < getType.NumField(); i++ {
		if getType.Field(i).Name == fld {
			value := getValue.Field(i).Interface()
			return value
		}
	}
	return nil
}

func getParsedQueryByKeyword(q string) map[string]string {
	//query
	//select id,name,age from indexName where name:ram,age:>40,+age:<=50,startDt>2022-01-01T01:01:00Z facets name limit 1, 10

	indexs := global.RegExkeyWords.FindAllStringIndex(strings.ToLower(q), -1)
	kw := make(map[string]string)
	curIdx := 0
	curKey := ""
	maxLen := len(q)
	for _, r := range indexs {
		//fmt.Println("curidex", r)
		if curIdx > 0 {
			kw[strings.ToLower(curKey[0:3])] = strings.Trim(q[curIdx:r[0]], " ")
		}
		curKey = q[r[0]:r[1]]
		curIdx = r[1]

	}
	if curIdx < maxLen {
		//fmt.Println("curidex", curIdx, maxLen)
		r := indexs[len(indexs)-1]

		curKey = q[r[0]:r[1]]
		kw[strings.ToLower(curKey[0:3])] = strings.Trim(q[r[1]:], " ")
	}
	//fmt.Println("query parser", kw)
	return kw
}
