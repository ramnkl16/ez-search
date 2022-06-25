package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/auth"
	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/ezcsv"
	"github.com/ramnkl16/ez-search/ezsearch"

	//"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

var (
	SearchController searchControllerInteface = &searchController{}
)

type searchControllerInteface interface {
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type searchController struct{}

func (ctrl *searchController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/search", ctrl.Search)
	rout.POST("/api/addorupdate", ctrl.AddOrUpdateIndex)
	rout.DELETE("/api/search", ctrl.DeleteIndex)
	rout.GET("/api/getindexes", ctrl.GetAllIndexes)
	rout.GET("/api/getfields", ctrl.GetFields)
	rout.POST("/api/createschema", ctrl.createIndexSchema)
	rout.GET("/api/getschema", ctrl.getIndexSchema)
	rout.POST("/api/generateSchema", ctrl.GenerateIndexSchema)

}

// bleve indexes
// getsearch  godoc
// @Summary get search result
// @Description get search result sample query [select id,name,age from indexName where name:ram,age:>40,+age:<=50,startDt>"2022-01-01T01:01:00Z" facets name limit 1, 10]
// @Description fetch record first 10 records with matching codition and shows sepecified fields
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Param  query body ezsearch.SearchRequestQuery true "look search"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/search [post]
func (ctrl *searchController) Search(ctx *gin.Context) {
	//fmt.Println("reached controlser|search")
	var m ezsearch.SearchRequestQuery

	if err := ctx.ShouldBindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	//searchResult, err := ezsearch.PostSearchResult(m)
	sr, err := ezsearch.PostSearchResult(m.Query)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	//searchResult, err := ezsearch.MergeSearchResult(sr)
	// if err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }

	// fmt.Println("search controller")
	// fmt.Println(m)
	if sr == nil {
		ctx.JSON(404, rest_errors.NewNotFoundError(m.Query))
		return
	}
	ctx.JSON(http.StatusCreated, sr)
}

// addorupdate  godoc
// @Summary Add or Update Index Document
// @Description AddorUpdate index
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Param  indexName query string true "name of the index you can also provide pattern like indexName{2006-01-02}-->indexName{yyyy-MM-dd}"
// @Param  docId query string false "document id"
// @Param  indexTranDate query string false "index name determind using this date when index name pattern date format must be {yyyy-MM-dd}"
// @Param  reqModel body interface{} true "document to be index"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/addorupdate [post]
func (ctrl *searchController) AddOrUpdateIndex(ctx *gin.Context) {

	indexName := ctx.Query("indexName")
	pindexName := indexName
	//fmt.Println("controlller|addorUpdate|", pindexName)
	//abstractimpl.IndexBasePath = ns
	docId := ctx.Query("docId")
	if len(docId) == 0 {
		docId = uid_utils.GetUid("uk", true)
	}
	ns := auth.GetNamespace(ctx.Request)
	if !strings.Contains(indexName, "/") {

		indexName = fmt.Sprintf("%s%c%s", ns, '/', indexName)
	}

	//fmt.Println("addorupdate", indexName, docId)
	//matches := m.FindAllString("index{2006-09-11}", -1)
	//dtVal := time.Now().UTC().Format(dtFormat)
	dt := ctx.Query("indexTranDate")
	if len(dt) == 0 && strings.Contains(indexName, "{") {
		format := common.ExtractDateFormatFromIndex(indexName)
		dt = time.Now().Format(format)
	}

	indexName, errMsg := common.GetPatternIndexName(indexName, dt)
	if len(errMsg) > 0 {
		ctx.String(http.StatusBadRequest, errMsg)
		return
	}
	//fmt.Printf("%q\n", global.RegexParseDate.FindAllSubmatch([]byte(`index{2006-09-11}`), -1))
	fmt.Println("ctrl|after pattern|indexName", indexName)
	var err error
	i, err := ezsearch.GetIndex(indexName)
	if i == nil || err != nil {
		//try to create new index

		//fmt.Println("Failed while update index|AddOrUpdateIndex", indexName, docId, saveErr)
		i, err = abstractimpl.BuilddynamicSchema(pindexName, indexName)
		ctx.JSON(http.StatusInternalServerError, err)
		return

	}

	// if len(docId) == 0 {
	// 	saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("unable to AddOrUpdateIndex [%s] for id[%s]. Please provide correct doc id", indexName, docId))
	// 	ctx.JSON(saveErr.Status(), saveErr)
	// 	return
	// }
	var m interface{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("AddOrUpdateIndex|invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//fmt.Println("docId and model", indexName, docId, m)
	err = i.Index(docId, m)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.String(http.StatusCreated, "success")
}

// delete  godoc
// @Summary Delete index document
// @Description Delete index
// @Tags bleve indexes
// @Produce  json
// @Param  indexName query string true "name of the index"
// @Param  docId query string true "document id"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/search [delete]
func (ctrl *searchController) DeleteIndex(ctx *gin.Context) {
	indexName := ctx.Query("indexName")
	ns := auth.GetNamespace(ctx.Request)
	if !strings.Contains(indexName, "/") {

		indexName = fmt.Sprintf("%s%c%s", ns, '/', indexName)
	}

	docId := ctx.Query("docId")
	var err error
	i, err := ezsearch.GetIndex(indexName)
	if i == nil {
		saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("index [%s] is not created. Please try after the index first", indexName))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(docId) == 0 {
		saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("unable to add or update index [%s] for id[%s]. Please provide correct doc id", indexName, docId))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	err = i.Delete(docId)
	//fmt.Println("delete account testing", err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.String(http.StatusCreated, "deleted")
}

// GetAllIndexes  godoc
// @Summary Get All indexes
// @Description Get all indexes
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Success 200 {object} []string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/getindexes [get]
func (ctrl *searchController) GetAllIndexes(ctx *gin.Context) {
	list := make([]string, 0)
	ns := auth.GetNamespace(ctx.Request)
	isPlatform := false
	if ns == "platform" {
		isPlatform = true
	}
	for k := range common.GetAllIndexes() {
		if isPlatform {
			list = append(list, k) //if users from platform should see all indexes otherwise only specific folder.
		} else if strings.HasPrefix(k, ns) {
			list = append(list, k)
		}
	}
	ctx.JSON(http.StatusCreated, list)
}

// GetFields  godoc
// @Summary GetFields
// @Description Get index fields
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Param  indexName query string true "name of the index"
// @Success 200 {object} []string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/getfields [get]
func (ctrl *searchController) GetFields(ctx *gin.Context) {
	indexName := ctx.Query("indexName")
	//fmt.Println("controllers|getfields", indexName)
	ctx.JSON(http.StatusCreated, ezsearch.GetFields(indexName))
}

// GetFields  godoc
// @Summary create schema
// @Description create schema possible values for type[bool|text|date|numeric|geopoint]
// @Description golang follows date pattern like indexname{2006-01-02} indexname{2006-01-02} which is equal to {yyyy-MM-dd}
// @Description indexname{2006-01-02}-->creates a index every day
// @Description indexname{2006-01}-->creates a index every month
// @Description indexname{2006}-->creates a index every year
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Param  indexName query string true "name of the index, you can also provide index date pattern like indexname{2006-01-02}"
// @Param  fieldDef body []ezsearch.BleveFieldDef true "field definition"
// @Success 200 {object} []string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/createschema [post]
func (ctrl *searchController) createIndexSchema(ctx *gin.Context) {
	indexName := ctx.Query("indexName")
	logger.Debug("createIndexSchema|controller", zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})

	var m []common.BleveFieldDef

	if err := ctx.ShouldBindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("CreateIndexSchema|invalid json body")
		logger.Error("", saveErr)
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ns := auth.GetNamespace(ctx.Request)
	logger.Debug("json", zapcore.Field{String: fmt.Sprintf("%v", m), Key: "p1", Type: zapcore.StringType})
	err := abstractimpl.BuildIndexSchema(indexName, m, ns)
	if err == nil {
		ctx.JSON(http.StatusCreated, "Success")
	} else {
		logger.Error("", err, zapcore.Field{Integer: int64(err.Status()), Key: "p1", Type: zapcore.StringType})
		ctx.JSON(err.Status(), err.Message())

	}

}

// GetFields  godoc
// @Summary get schema
// @Description get schema possible values for type[bool|text|date|numeric|geopoint]
// @Description golang follows date pattern like indexname{2006-01-02} indexname{2006-01-02} which is equal to {yyyy-MM-dd}
// @Description indexname{2006-01-02}-->gets a index every day
// @Description indexname{2006-01}-->gets a index every month
// @Description indexname{2006}-->gets a index every year
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Param  indexName query string true "name of the index, you can also provide index date pattern like indexname{2006-01-02}"
// @Param  fieldDef body []ezsearch.BleveFieldDef true "field definition"
// @Success 200 {object} []string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/getschema [get]
func (ctrl *searchController) getIndexSchema(ctx *gin.Context) {
	indexName := ctx.Query("indexName")
	str, _ := url.QueryUnescape(indexName)
	//fmt.Println("url decode", str)
	logger.Debug("getIndexSchema|controller", zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
	key := fmt.Sprintf("%s.schema", str)

	schemaByte, err := coredb.GetKey(key)
	if err != nil {
		//errStr := fmt.Sprintf(`%s schema is not found in core db. Please try after get schema first. \n%s\n`, indexName, err.Error())
		schemaErr := rest_errors.NewBadRequestError("%s schema is not found in core db. Please try after create schema first")
		ctx.JSON(schemaErr.Status(), schemaErr)
		return
	}
	if len(schemaByte) == 0 {
		schemaErr := rest_errors.NewBadRequestError("%s schema is not found in core db. Please try after create schema first")
		ctx.JSON(schemaErr.Status(), schemaErr)
		return
	}
	//fmt.Println("getschema", string(schemaByte))
	ctx.String(http.StatusCreated, string(schemaByte))

}

// GetFields  godoc
// @Summary get schema
// @Description get schema possible values for type[bool|text|date|numeric|geopoint]
// @Description golang follows date pattern like indexname{2006-01-02} indexname{2006-01-02} which is equal to {yyyy-MM-dd}
// @Description indexname{2006-01-02}-->gets a index every day
// @Description indexname{2006-01}-->gets a index every month
// @Description indexname{2006}-->gets a index every year
// @Tags bleve indexes
// @Accept  json
// @Produce  json
// @Param  indexName query string true "name of the index, you can also provide index date pattern like indexname{2006-01-02}"
// @Param  fieldDef body []ezsearch.BleveFieldDef true "field definition"
// @Success 200 {object} []string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/generateSchema [post]
type schemaString struct {
	Columns     string `json:"columns"`
	IsShortName bool   `json:"isShortName"`
}

func (ctrl *searchController) GenerateIndexSchema(ctx *gin.Context) {

	var m schemaString

	if err := ctx.ShouldBindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	list := ezcsv.GenerateIndexSchema(m.Columns, m.IsShortName)
	if len(list) == 0 {
		schemaErr := rest_errors.NewBadRequestError("Failed while generate schema")
		ctx.JSON(schemaErr.Status(), schemaErr)
		return
	}
	ctx.JSON(http.StatusCreated, list)
	//generate column list
	cols := make([]string, 0)
	for _, v := range list {
		cols = append(cols, v.Name)
	}
	ctx.String(http.StatusAccepted, strings.Join(cols, ","))

}
