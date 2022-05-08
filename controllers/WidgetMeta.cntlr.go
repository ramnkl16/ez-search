package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/auth"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
)

var (
	WidgetMetaController widgetmetaControllerInteface = &widgetmetaController{}
)

type widgetmetaControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

//Data source schema to persist bleve search query

type widgetmetaController struct{}

func (ctrl *widgetmetaController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/widgetmeta", ctrl.Create)
	rout.PUT("/api/widgetmeta", ctrl.Update)
	rout.DELETE("/api/widgetmeta/:id", ctrl.Delete)
	rout.GET("/api/widgetmeta/:id", ctrl.Get)
	rout.GET("/api/widgetmetas/search", ctrl.Search)
	//rout.GET("/web-ui/widgetmetas/search", ctrl.Search)
}

//Create new WidgetMeta

// Create  godoc
// @Summary Create WidgetMeta
// @Description create WidgetMeta
// @Tags WidgetMeta
// @Accept  json
// @Produce  json
// @Param  address body models.WidgetMeta true "create WidgetMeta"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/widgetmeta [post]
func (ctrl *widgetmetaController) Create(ctx *gin.Context) {

	var wi models.WidgetMeta

	if err := ctx.ShouldBindJSON(&wi); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.WidgetMetaService.Create(wi); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created WidgetMeta "})
}

// Update  godoc
// @Summary Update WidgetMeta
// @Description Update by json WidgetMeta
// @Tags WidgetMeta
// @Accept  json
// @Produce  json
// @Param  WidgetMeta body models.WidgetMeta true "Update WidgetMeta"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/widgetmeta [put]
func (ctrl *widgetmetaController) Update(ctx *gin.Context) {

	var wi models.WidgetMeta

	if err := ctx.ShouldBindJSON(&wi); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.WidgetMetaService.Update(wi); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "WidgetMeta updated"})
}

// Search godoc
// @Summary Get all WidgetMeta
// @Description get top 100 records WidgetMeta
// @Tags WidgetMeta
// @Produce  json
// @Param  start query int true "starting row"
// @Param  limit query int true "no of row limit"
// @Success 200 {array} models.WidgetMeta
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/widgetmetas/search [get]
//getting all  WidgetMeta
func (ctrl *widgetmetaController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	limit := ctx.Query("limit")

	results, err := services.WidgetMetaService.Search(start, limit)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ns := auth.GetNamespace(ctx.Request)
	filteredList := make([]models.WidgetMeta, 0)
	if ns != "platform" {
		for _, item := range results {
			if item.Division == ns {
				filteredList = append(filteredList, item)
			}
		}
		ctx.JSON(http.StatusOK, filteredList) //	return filteredList, nil
		return
	}
	ctx.JSON(http.StatusOK, results)
}

// Get godoc
// @Summary Get all WidgetMeta
// @Description Get WidgetMeta by id
// @Tags WidgetMeta
// @Produce  json
// @Param  id path string true "WidgetMeta ID"
// @Success 200 {array} models.WidgetMeta
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/widgetmeta/{id} [get]
//getting a record  WidgetMeta
func (ctrl *widgetmetaController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	wi, err := services.WidgetMetaService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, wi)
}

// Delete  godoc
// @Summary Get all WidgetMeta
// @Description delete WidgetMeta
// @Tags WidgetMeta
// @Produce  json
// @Param  id path string true "WidgetMeta ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/widgetmeta/{id} [delete]
//delete  WidgetMeta by id
func (wic *widgetmetaController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.WidgetMetaService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "WidgetMeta deleted"})
}
