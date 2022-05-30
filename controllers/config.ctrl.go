package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/rest_errors"
)

var (
	ConfigController configControllerInteface = &configController{}
)

type configControllerInteface interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type configController struct{}

func (ctrl *configController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/config", ctrl.Create)
	rout.DELETE("/api/config/:key", ctrl.Delete)
	rout.GET("/api/config/:key", ctrl.Get)

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

func (ctrl *configController) Create(ctx *gin.Context) {
	//fmt.Println("reached controlser|Create")
	var m regConfigJson

	if err := ctx.BindJSON(&m); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body|config.Create")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	err := coredb.AddKey(m.Key, []byte(m.Value))

	if err != nil {
		saveErr := rest_errors.NewInternalServerError("Failed while add in Coredb|config.create", err)
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, "Register successfully.")
}

// Get godoc
// @Summary Get all config
// @Description Get config by key
// @Tags config
// @Produce  json
// @Param  id path string true "config key"
// @Success 200 {array} regConfigJson
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/config/{key} [get]
//getting a record  config
func (ctrl *configController) Get(ctx *gin.Context) {
	fmt.Println("config|get|controller")
	id := ctx.Param("key")
	ev, err := coredb.GetKey(id)
	if err != nil {
		restErr := rest_errors.NewInternalServerError("Failed while get key from core", err)
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	ctx.String(http.StatusOK, string(ev))
}

// Delete  godoc
// @Summary Get all config
// @Description delete config
// @Tags config
// @Produce  json
// @Param  id path string true "config ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/config/{id} [delete]
//delete  config by id
func (evc *configController) Delete(ctx *gin.Context) {
	id := ctx.Param("key")
	err := coredb.Delete(id)
	if err != nil {
		restErr := rest_errors.NewInternalServerError("Failed while deleted from coredb", err)
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	ctx.String(http.StatusOK, fmt.Sprintf("Deleted %s", id))
}
