package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
)

var (
	EventQueueController eventqueueControllerInteface = &eventqueueController{}
)

type eventqueueControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type eventqueueController struct{}

func (ctrl *eventqueueController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/eventqueue", ctrl.Create)
	rout.PUT("/api/eventqueue", ctrl.Update)
	rout.DELETE("/api/eventqueue/:id", ctrl.Delete)
	rout.GET("/api/eventqueue/:id", ctrl.Get)
	rout.GET("/api/eventqueues/search", ctrl.Search)
	//rout.GET("/web-ui/eventqueues/search", ctrl.Search)
}

//Create new EventQueue

// Create  godoc
// @Summary Create EventQueue
// @Description create EventQueue
// @Tags EventQueue
// @Accept  json
// @Produce  json
// @Param  address body models.EventQueue true "create EventQueue"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/eventqueue [post]
func (ctrl *eventqueueController) Create(ctx *gin.Context) {

	var ev models.EventQueue

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.EventQueueCustomService.CreateWithIndex(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created EventQueue "})
}

// Update  godoc
// @Summary Update EventQueue
// @Description Update by json EventQueue
// @Tags EventQueue
// @Accept  json
// @Produce  json
// @Param  EventQueue body models.EventQueue true "Update EventQueue"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/eventqueue [put]
func (ctrl *eventqueueController) Update(ctx *gin.Context) {

	var ev models.EventQueue

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.EventQueueCustomService.UpdateWithIndex(ev); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "EventQueue updated"})
}

// Search godoc
// @Summary Get all EventQueue
// @Description get top 100 records EventQueue
// @Tags EventQueue
// @Produce  json
// @Param  start query int false "Starting row 0-->first record"
// @Param  limit query int false "Row limit"
// @Success 200 {array} models.EventQueue
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/eventqueues/search [get]
//getting all  EventQueue
func (ctrl *eventqueueController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	if len(start) == 0 {
		start = "0"
	}
	limit := ctx.Query("limit")
	if len(limit) == 0 {
		limit = "50"
	}
	query := fmt.Sprintf("select * from %s limit %s,%s", abstractimpl.EventQueueTable, start, limit)
	//fmt.Println("evenquery query", query)
	results, err := services.EventQueueService.Search(query)
	if err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, results)
}

// Get godoc
// @Summary Get all EventQueue
// @Description Get EventQueue by id
// @Tags EventQueue
// @Produce  json
// @Param  id path string true "EventQueue ID"
// @Success 200 {array} models.EventQueue
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/eventqueue/{id} [get]
//getting a record  EventQueue
func (ctrl *eventqueueController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	ev, err := services.EventQueueService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, ev)
}

// Delete  godoc
// @Summary Get all EventQueue
// @Description delete EventQueue
// @Tags EventQueue
// @Produce  json
// @Param  id path string true "EventQueue ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/eventqueue/{id} [delete]
//delete  EventQueue by id
func (evc *eventqueueController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.EventQueueService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "EventQueue deleted"})
}
