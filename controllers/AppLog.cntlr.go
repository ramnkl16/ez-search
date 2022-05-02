package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
)

var (
	AppLogController appLogControllerInteface = &appLogController{}
)

type appLogControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type appLogController struct{}

func (ctrl *appLogController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/appLog", ctrl.Create)
	rout.PUT("/api/appLog", ctrl.Update)
	rout.DELETE("/api/appLog/:id", ctrl.Delete)
	rout.GET("/api/appLog/:id", ctrl.Get)
	rout.GET("/api/appLogs/search", ctrl.Search)
	//rout.GET("/web-ui/appLogs/search", ctrl.Search)
}

func (ctrl *appLogController) Create(ctx *gin.Context) {

	var ap models.AppLog

	if err := ctx.ShouldBindJSON(&ap); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	// if err := services.LogsCustomJobService.CreateAppLog(ap); err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created AppLog "})
}

func (ctrl *appLogController) Update(ctx *gin.Context) {

	var ap models.AppLog

	if err := ctx.ShouldBindJSON(&ap); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	// if err := services.AppLogService.Update(ap); err != nil {
	// 	ctx.JSON(err.Status(), err)
	// }
	ctx.JSON(http.StatusOK, map[string]string{"status": "AppLog updated"})
}

func (ctrl *appLogController) Search(ctx *gin.Context) {
	// start := ctx.Query("start")
	// limit := ctx.Query("limit")

	// results, err := services.AppLogService.Search(start, limit)
	// if err != nil {
	// 	ctx.JSON(err.Status(), err)
	// }
	// ctx.JSON(http.StatusOK, results)
}

func (ctrl *appLogController) Get(ctx *gin.Context) {
	// id := ctx.Param("id")
	// ap, err := services.AppLogService.Get(id)
	// if err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }
	// ctx.JSON(http.StatusOK, ap)
}

func (apc *appLogController) Delete(ctx *gin.Context) {
	//id := ctx.Param("id")
	// err := services.AppLogService.Delete(id)
	// if err != nil {
	// 	ctx.JSON(err.Status(), err)
	// 	return
	// }
	// ctx.JSON(http.StatusOK, map[string]string{"status": "AppLog deleted"})
}
