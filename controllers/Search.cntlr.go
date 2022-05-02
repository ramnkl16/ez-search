package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/ezsearch"
	"github.com/ramnkl16/ez-search/global"
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
	fmt.Println("reached controlser|search")
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
	docId := ctx.Query("docId")
	if len(docId) == 0 {
		docId = uid_utils.GetUid("uk", true)
	}

	//fmt.Println("addorupdate", indexName, docId)
	//matches := m.FindAllString("index{2006-09-11}", -1)
	//dtVal := time.Now().UTC().Format(dtFormat)
	dt := ctx.Query("indexTranDate")
	indexName, errMsg := common.GetPatternIndexName(indexName, dt)
	if len(errMsg) > 0 {
		ctx.String(http.StatusBadRequest, errMsg)
		return
	}
	fmt.Printf("%q\n", global.RegexParseDate.FindAllSubmatch([]byte(`index{2006-09-11}`), -1))

	var err error
	i, err := ezsearch.GetIndex(indexName)
	if i == nil || err != nil {
		//try to create new index

		fdefs, err := ezsearch.GetBleveTableschema(pindexName)
		if err != nil {
			saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("index [%s] schema was not created", indexName))
			logger.Error(fmt.Sprintf("Failed|AddOrUpdateIndex id=%s, [%s]", indexName, docId), saveErr)
		}
		if fdefs == nil || len(fdefs) == 0 {
			fdefs = append(fdefs, common.BleveFieldDef{Name: "timestamp", Type: "date"})
		}
		err = abstractimpl.BuildIndexSchema(indexName, fdefs, false)
		i, err = ezsearch.GetIndex(indexName)
		if err != nil {
			saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("index [%s] is not created. Please try after the index first", indexName))
			logger.Error(fmt.Sprintf("Failed|AddOrUpdateIndex id=%s, [%s]", indexName, docId), saveErr)
			//fmt.Println("Failed while update index|AddOrUpdateIndex", indexName, docId, saveErr)
		}

		//return
	}

	if len(docId) == 0 {
		saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("unable to AddOrUpdateIndex [%s] for id[%s]. Please provide correct doc id", indexName, docId))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	var m interface{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("AddOrUpdateIndex|invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//fmt.Println("docId and model", docId, m)
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
	for k, _ := range common.GetAllIndexes() {
		list = append(list, k)
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
	logger.Debug("json", zapcore.Field{String: fmt.Sprintf("%v", m), Key: "p1", Type: zapcore.StringType})
	err := abstractimpl.BuildIndexSchema(indexName, m, false)
	if err == nil {
		ctx.JSON(http.StatusCreated, "Success")
	} else {
		logger.Error("", err, zapcore.Field{Integer: int64(err.Status()), Key: "p1", Type: zapcore.StringType})
		ctx.JSON(err.Status(), err.Message())
	}

}