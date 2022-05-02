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
	NamespaceController namespaceControllerInteface = &namespaceController{}
)

type namespaceControllerInteface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	New(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type namespaceController struct{}

func (ctrl *namespaceController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/namespace", ctrl.Create)
	rout.PUT("/api/namespace", ctrl.Update)
	rout.DELETE("/api/namespace/:id", ctrl.Delete)
	rout.GET("/api/namespace/:id", ctrl.Get)
	rout.GET("/api/namespaces/search", ctrl.Search)
	rout.POST("/api/namespacenew", ctrl.New)
	//rout.GET("/web-ui/namespaces/search", ctrl.Search)
}

//Create new Namespace
// Create  godoc
// @Summary Create Namespace
// @Description create Namespace
// @Tags Namespace
// @Accept  json
// @Produce  json
// @Param  address body models.Namespace true "create Namespace"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/namespacenew [post]
func (ctrl *namespaceController) New(ctx *gin.Context) {

	var ev models.NamespaceParam

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.NamespaceService.New(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "New Namespace "})
}

//new new Namespace
// Create  godoc
// @Summary new Namespace
// @Description create Namespace
// @Tags Namespace
// @Accept  json
// @Produce  json
// @Param  address body models.NamespaceParam true "create NamespaceParam"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/namespace [post]
func (ctrl *namespaceController) Create(ctx *gin.Context) {

	var ev models.Namespace

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if err := services.NamespaceService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusCreated, map[string]string{"status": "created Namespace "})
}

// Update  godoc
// @Summary Update Namespace
// @Description Update by json Namespace
// @Tags Namespace
// @Accept  json
// @Produce  json
// @Param  Namespace body models.Namespace true "Update Namespace"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/namespace [put]
func (ctrl *namespaceController) Update(ctx *gin.Context) {

	var ev models.Namespace

	if err := ctx.ShouldBindJSON(&ev); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := services.NamespaceService.Save(ev); err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "Namespace updated"})
}

// Search godoc
// @Summary Get all Namespace
// @Description get top 100 records Namespace
// @Tags Namespace
// @Produce  json
// @Param  start query int false "Starting row 0-->first record"
// @Param  limit query int false "Row limit"
// @Success 200 {array} models.Namespace
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/namespaces/search [get]
//getting all  Namespace
func (ctrl *namespaceController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	if len(start) == 0 {
		start = "0"
	}
	limit := ctx.Query("limit")
	if len(limit) == 0 {
		limit = "50"
	}
	query := fmt.Sprintf("select * from %s limit %s,%s", abstractimpl.NamespaceTable, start, limit)
	//fmt.Println("evenquery query", query)
	results, err := services.NamespaceService.Search(query)
	if err != nil {
		ctx.JSON(err.Status(), err)
	}
	ctx.JSON(http.StatusOK, results)
}

// Get godoc
// @Summary Get all Namespace
// @Description Get Namespace by id
// @Tags Namespace
// @Produce  json
// @Param  id path string true "Namespace ID"
// @Success 200 {array} models.Namespace
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/namespace/{id} [get]
//getting a record  Namespace
func (ctrl *namespaceController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	ev, err := services.NamespaceService.Get(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, ev)
}

// Delete  godoc
// @Summary Get all Namespace
// @Description delete Namespace
// @Tags Namespace
// @Produce  json
// @Param  id path string true "Namespace ID"
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/namespace/{id} [delete]
//delete  Namespace by id
func (evc *namespaceController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.NamespaceService.Delete(id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "Namespace deleted"})
}
