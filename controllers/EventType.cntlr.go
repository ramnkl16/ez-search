package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
)

var (
	EventTypeController eventtypeControllerInteface = &eventtypeController{}
)

type eventtypeControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type eventtypeController struct{}

func (ctrl *eventtypeController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/eventtype", ctrl.Create)
	rout.PUT("/api/eventtype", ctrl.Update)
	rout.DELETE("/api/eventtype/:id", ctrl.Delete)
	rout.GET("/api/eventtype/:id", ctrl.Get)
	rout.GET("/api/eventtypes/search", ctrl.Search)
	//rout.GET("/web-ui/eventtypes/search", ctrl.Search)
}

func (ctrl *eventtypeController) Create(ctx *gin.Context) {

	var ev models.EventType

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	// if err := services.EventTypeService.Create(ev); err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created EventType "})
}

func (ctrl *eventtypeController) Update(ctx *gin.Context) {

	var ev models.EventType

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	// if err := services.EventTypeService.Update(ev); err != nil {
	// 	ctx.JSON(err.Status(), err)
	// }
	ctx.JSON(http.StatusOK, map[string]string{"status": "EventType updated"})
}

func (ctrl *eventtypeController) Search(ctx *gin.Context) {
	// start := ctx.Query("start")
	// limit := ctx.Query("limit")

	// // results, err := services.EventTypeService.Search(start, limit)
	// // if err != nil {
	// // 	ctx.JSON(err.Status(), err)
	// // }
	// ctx.JSON(http.StatusOK, results)

}

func (ctrl *eventtypeController) Get(ctx *gin.Context) {
	// id := ctx.Param("id")
	// ev, err := services.EventTypeService.Get(id)
	// if err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }
	// ctx.JSON(http.StatusOK, ev)
}

func (evc *eventtypeController) Delete(ctx *gin.Context) {
	// id := ctx.Param("id")
	// err := services.EventTypeService.Delete(id)
	// if err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }
	// ctx.JSON(http.StatusOK, map[string]string{"status": "EventType deleted"})
}
