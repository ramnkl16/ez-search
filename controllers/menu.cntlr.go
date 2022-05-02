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
	MenuController menuControllerInteface = &menuController{}
)

type menuControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type menuController struct{}

func (ctrl *menuController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/menu", ctrl.Create)
	rout.PUT("/api/menu", ctrl.Update)
	rout.DELETE("/api/menu/:id", ctrl.Delete)
	rout.GET("/api/menu/:id", ctrl.Get)
	rout.GET("/api/menus/search", ctrl.Search)
	//rout.GET("/web-ui/menus/search", ctrl.Search)
}

//Create new menu

// Create  godoc
// @Summary Create Menu
// @Description create Menu
// @Tags Menu
// @Accept  json
// @Produce  json
// @Param  address body models.Menu true "create Menu"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/menu [post]
func (ctrl *menuController) Create(ctx *gin.Context) {

	var ev models.Menu

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.MenuService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created Menu "})
}

// Update  godoc
// @Summary Update Menu
// @Description Update by json Menu
// @Tags Menu
// @Accept  json
// @Produce  json
// @Param  Menu body models.Menu true "Update Menu"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/menu [put]
func (ctrl *menuController) Update(ctx *gin.Context) {

	var ev models.Menu

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.MenuService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "Menu updated"})
}

// Search godoc
// @Summary Get all Menu
// @Description get top 100 records Menu
// @Tags Menu
// @Produce  json
// @Param  start query int false "Starting row 0-->first record"
// @Param  limit query int false "Row limit"
// @Success 200 {array} models.Menu
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/menus/search [get]
//getting all  Menu
func (ctrl *menuController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	if len(start) == 0 {
		start = "0"
	}
	limit := ctx.Query("limit")
	if len(limit) == 0 {
		limit = "50"
	}
	query := fmt.Sprintf("select * from %s limit %s,%s", abstractimpl.MenuTable, start, limit)
	//fmt.Println("evenquery query", query)
	results, err := services.MenuService.Search(query)
	if err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, results)
}

// Get godoc
// @Summary Get all Menu
// @Description Get Menu by id
// @Tags Menu
// @Produce  json
// @Param  id path string true "menu ID"
// @Success 200 {array} models.menu
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/menu/{id} [get]
//getting a record  Menu
func (ctrl *menuController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	ev, err := services.MenuService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, ev)
}

// Delete  godoc
// @Summary Get all menu
// @Description delete menu
// @Tags menu
// @Produce  json
// @Param  id path string true "menu ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/menu/{id} [delete]
//delete  menu by id
func (evc *menuController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.MenuService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "menu deleted"})
}
