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
	UserGroupController userGroupControllerInteface = &userGroupController{}
)

type userGroupControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type userGroupController struct{}

func (ctrl *userGroupController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/userGroup", ctrl.Create)
	rout.PUT("/api/userGroup", ctrl.Update)
	rout.DELETE("/api/userGroup/:id", ctrl.Delete)
	rout.GET("/api/userGroup/:id", ctrl.Get)
	rout.GET("/api/userGroups/search", ctrl.Search)
	//rout.GET("/web-ui/userGroups/search", ctrl.Search)
}

//Create new UserGroup

// Create  godoc
// @Summary Create UserGroup
// @Description create UserGroup
// @Tags UserGroup
// @Accept  json
// @Produce  json
// @Param  address body models.UserGroup true "create UserGroup"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userGroup [post]
func (ctrl *userGroupController) Create(ctx *gin.Context) {

	var ev models.UserGroup

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.UserGroupService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created userGroup "})
}

// Update  godoc
// @Summary Update UserGroup
// @Description Update by json UserGroup
// @Tags UserGroup
// @Accept  json
// @Produce  json
// @Param  UserGroup body models.UserGroup true "Update UserGroup"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userGroup [put]
func (ctrl *userGroupController) Update(ctx *gin.Context) {

	var ev models.UserGroup

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.UserGroupService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "UserGroup updated"})
}

// Search godoc
// @Summary Get all UserGroup
// @Description get top 100 records UserGroup
// @Tags UserGroup
// @Produce  json
// @Param  start query int false "Starting row 0-->first record"
// @Param  limit query int false "Row limit"
// @Success 200 {array} models.UserGroup
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userGroups/search [get]
//getting all  UserGroup
func (ctrl *userGroupController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	if len(start) == 0 {
		start = "0"
	}
	limit := ctx.Query("limit")
	if len(limit) == 0 {
		limit = "50"
	}
	query := fmt.Sprintf("select * from %s limit %s,%s", abstractimpl.UserGroupTable, start, limit)
	//fmt.Println("evenquery query", query)
	results, err := services.UserGroupService.Search(query)
	if err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, results)
}

// Get godoc
// @Summary Get all UserGroup
// @Description Get UserGroup by id
// @Tags UserGroup
// @Produce  json
// @Param  id path string true "UserGroup ID"
// @Success 200 {array} models.UserGroup
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userGroup/{id} [get]
//getting a record  UserGroup
func (ctrl *userGroupController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	ev, err := services.UserGroupService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, ev)
}

// Delete  godoc
// @Summary Get all UserGroup
// @Description delete UserGroup
// @Tags UserGroup
// @Produce  json
// @Param  id path string true "UserGroup ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/userGroup/{id} [delete]
//delete  UserGroup by id
func (evc *userGroupController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.UserGroupService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "UserGroup deleted"})
}
