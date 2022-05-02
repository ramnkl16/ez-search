package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
)

var (
	EventQueueHistoryController eventqueuehistoryControllerInteface = &eventqueuehistoryController{}
)

type eventqueuehistoryControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type eventqueuehistoryController struct{}

func (ctrl *eventqueuehistoryController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/eventqueuehistory", ctrl.Create)
	rout.PUT("/api/eventqueuehistory", ctrl.Update)
	rout.DELETE("/api/eventqueuehistory/:id", ctrl.Delete)
	rout.GET("/api/eventqueuehistory/:id", ctrl.Get)
	rout.GET("/api/eventqueuehistories/search", ctrl.Search)
	//rout.GET("/web-ui/eventqueuehistories/search", ctrl.Search)
}

func (ctrl *eventqueuehistoryController) Create(ctx *gin.Context) {

	var ev models.EventQueueHistory

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.EventQueueHistoryService.Create(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created EventQueueHistory "})
}

func (ctrl *eventqueuehistoryController) Update(ctx *gin.Context) {

	var ev models.EventQueueHistory

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.EventQueueHistoryService.Update(ev); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "EventQueueHistory updated"})
}

func (ctrl *eventqueuehistoryController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	limit := ctx.Query("limit")

	results, err := services.EventQueueHistoryService.Search(start, limit)
	if err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, results)
}

func (ctrl *eventqueuehistoryController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	ev, err := services.EventQueueHistoryService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, ev)
}

func (evc *eventqueuehistoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.EventQueueHistoryService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "EventQueueHistory deleted"})
}
