package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/rest_errors"
)

var (
	CoredbController coredbControllerInteface = &coredbController{}
)

type coredbControllerInteface interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetValues(ctx *gin.Context)
	GetKeys(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type coredbController struct{}

func (ctrl *coredbController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/coredb", ctrl.Create)
	rout.DELETE("/api/coredb/:key", ctrl.Delete)
	rout.GET("/api/coredb/:key", ctrl.Get)
	rout.GET("/api/coredb/getvalues", ctrl.GetValues)
	rout.GET("/api/coredb/getkeys", ctrl.GetKeys)

}

type regConfigJson struct {
	Value string `json:"value"`
	Key   string `json:"key"`
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

func (ctrl *coredbController) Create(ctx *gin.Context) {
	//fmt.Println("reached controlser|Create")
	var m regConfigJson

	if err := ctx.BindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body|coredb.Create")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	err := coredb.AddKey(coredb.Defaultbucket, m.Key, []byte(m.Value))

	if err != nil {
		saveErr := rest_errors.NewInternalServerError("Failed while add in Coredb|coredb.create", err)
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, "Register successfully.")
}

// Get godoc
// @Summary Get all coredb
// @Description Get coredb by key
// @Tags coredb
// @Produce  json
// @Param  id path string true "coredb key"
// @Success 200 {array} regConfigJson
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/coredb/{key} [get]
//getting a record  coredb
func (ctrl *coredbController) Get(ctx *gin.Context) {
	//fmt.Println("coredb|get|controller")
	id := ctx.Param("key")
	ev, err := coredb.GetValue(coredb.Defaultbucket, id)
	if err != nil {
		restErr := rest_errors.NewInternalServerError("Failed while get key from core", err)
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	ctx.String(http.StatusOK, string(ev))
}

// Delete  godoc
// @Summary Get all coredb
// @Description delete coredb
// @Tags coredb
// @Produce  json
// @Param  id path string true "coredb ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/coredb/{id} [delete]
//delete  coredb by id
func (evc *coredbController) Delete(ctx *gin.Context) {
	id := ctx.Param("key")
	err := coredb.Delete(coredb.Defaultbucket, id)
	if err != nil {
		restErr := rest_errors.NewInternalServerError("Failed while deleted from coredb", err)
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	ctx.String(http.StatusOK, fmt.Sprintf("Deleted %s", id))
}

// Get godoc
// @Summary Get all coredb
// @Description Get all values
// @Tags coredb
// @Produce  json
// @Param  id path string true "coredb key"
// @Success 200 {array} regcoredbJson
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/coredb/getvalues [get]
//getting a record  coredb
func (ctrl *coredbController) GetValues(ctx *gin.Context) {
	ev := coredb.GetValues(coredb.Defaultbucket)
	ctx.JSON(http.StatusOK, ev)
}

// Get godoc
// @Summary Get all coredb
// @Description Get all keys
// @Tags coredb
// @Produce  json
// @Param  id path string true "coredb key"
// @Success 200 {array} regcoredbJson
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/coredb/getvalues [get]
//getting a record  coredb
func (ctrl *coredbController) GetKeys(ctx *gin.Context) {
	ev := coredb.GetKeys(coredb.Defaultbucket)
	// if err != nil {
	// 	restErr := rest_errors.NewInternalServerError("Failed while get key from core", err)
	// 	ctx.JSON(restErr.Status(), restErr)
	// 	return
	// }
	ctx.JSON(http.StatusOK, ev)
}
