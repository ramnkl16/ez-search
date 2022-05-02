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
	UserMenuController userMenuControllerInteface = &userMenuController{}
)

type userMenuControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type userMenuController struct{}

func (ctrl *userMenuController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/userMenu", ctrl.Create)
	rout.PUT("/api/userMenu", ctrl.Update)
	rout.DELETE("/api/userMenu/:id", ctrl.Delete)
	rout.GET("/api/userMenu/:id", ctrl.Get)
	rout.GET("/api/userMenus/search", ctrl.Search)
	//rout.GET("/web-ui/userMenus/search", ctrl.Search)
}

//Create new userMenu

// Create  godoc
// @Summary Create UserMenu
// @Description create UserMenu
// @Tags UserMenu
// @Accept  json
// @Produce  json
// @Param  address body models.UserMenu true "create UserMenu"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userMenu [post]
func (ctrl *userMenuController) Create(ctx *gin.Context) {

	var ev models.UserMenu

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.UserMenuService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created userMenu "})
}

// Update  godoc
// @Summary Update UserMenu
// @Description Update by json UserMenu
// @Tags UserMenu
// @Accept  json
// @Produce  json
// @Param  UserMenu body models.UserMenu true "Update UserMenu"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userMenu [put]
func (ctrl *userMenuController) Update(ctx *gin.Context) {

	var ev models.UserMenu

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.UserMenuService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "UserMenu updated"})
}

// Search godoc
// @Summary Get all UserMenu
// @Description get top 100 records UserMenu
// @Tags UserMenu
// @Produce  json
// @Param  start query int false "Starting row 0-->first record"
// @Param  limit query int false "Row limit"
// @Success 200 {array} models.UserMenu
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userMenus/search [get]
//getting all  UserMenu
func (ctrl *userMenuController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	if len(start) == 0 {
		start = "0"
	}
	limit := ctx.Query("limit")
	if len(limit) == 0 {
		limit = "50"
	}
	query := fmt.Sprintf("select * from %s limit %s,%s", abstractimpl.UserMenuTable, start, limit)
	//fmt.Println("evenquery query", query)
	results, err := services.UserMenuService.Search(query)
	if err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, results)
}

// Get godoc
// @Summary Get all UserMenu
// @Description Get UserMenu by id
// @Tags UserMenu
// @Produce  json
// @Param  id path string true "UserMenu ID"
// @Success 200 {array} models.UserMenu
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userMenu/{id} [get]
//getting a record  UserMenu
func (ctrl *userMenuController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	ev, err := services.UserMenuService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, ev)
}

// Delete  godoc
// @Summary Get all UserMenu
// @Description delete UserMenu
// @Tags UserMenu
// @Produce  json
// @Param  id path string true "UserMenu ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userMenu/{id} [delete]
//delete  UserMenu by id
func (evc *userMenuController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.UserMenuService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "UserMenu deleted"})
}
