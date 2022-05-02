package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/auth"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
)

var (
	UserController userControllerInteface = &userController{}
)

type (
	LoginParam struct {
		EmailOrMobile string
		Password      string
		NsCode        string `json:"nsCode"`
	}
	userControllerInteface interface {
		Save(ctx *gin.Context)
		Get(ctx *gin.Context)
		Search(ctx *gin.Context)
		RegisterRouter(rout *gin.Engine)
		Login(rout *gin.Context)
		ChangePassword(rout *gin.Context)
		Delete(ctx *gin.Context)
		Logout(ctx *gin.Context)
	}
	changePasswordParam struct {
		EmailOrMobile string
		OldPassword   string
		NewPassword   string
	}
)

type userController struct{}

func (ctrl *userController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/user", ctrl.Save)
	rout.POST("/api/auth/login", ctrl.Login)
	rout.GET("/api/auth/logout", ctrl.Logout)
	rout.POST("/api/auth/changepassword", ctrl.ChangePassword)
	rout.GET("/api/user/:id", ctrl.Get)
	rout.GET("/api/users/search", ctrl.Search)
	rout.DELETE("/api/user/:id", ctrl.Delete)
}

//CustomSave  User

// Create/updte/delete  godoc
// @Summary Create/update/delete User
// @Description Save  User
// @Tags User
// @Accept  json
// @Produce  json
// @Param  address body models.User true "CustomSave models.User"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
//@Router /api/user [post]
func (ctrl *userController) Save(ctx *gin.Context) {

	var us models.User

	if err := ctx.ShouldBindJSON(&us); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body") //, rest_errors.InternalServerErrContactAdmin)

		handleFaultResponse(ctx, saveErr, rest_errors.InternalServerErrContactAdmin)
		return
	}
	res, err := services.UserService.Save(us)
	if err != nil {

		handleFaultResponse(ctx, err, rest_errors.InternalServerErrContactAdmin)
		return
	}
	handleSuccessResponse(ctx, http.StatusOK, map[string]string{"id": *res})
}

type loginResult struct {
	AuthToken   string `json:"authToken"`
	GroupId     string `json:"groupId"`
	NameSpaceId string `json:"namespaceID"`
}

//  auth login  godoc
// @Summary authenticate user
// @Description Save  User
// @Tags User
// @Accept  json
// @Produce  json
// @Param  user body LoginParam true "login LoginParam"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
//@Router /api/auth/login [post]
func (ctrl *userController) Login(ctx *gin.Context) {
	fmt.Println("login start")
	var lp LoginParam

	if err := ctx.ShouldBindJSON(&lp); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body|Login") //, rest_errors.InternalServerErrContactAdmin)
		// saveErr.SetErrorCode(rest_errors.InvalidUserIdOrPassword)
		// ctx.JSON(saveErr.Status(), saveErr)
		handleFaultResponse(ctx, saveErr, rest_errors.InvalidUserIdOrPassword)
		return
	}
	us, err := services.UserService.Login(lp.EmailOrMobile, lp.Password, lp.NsCode)
	if err != nil {
		//err.SetErrorCode(rest_errors.UnableToGenerateToken)
		handleFaultResponse(ctx, err, rest_errors.UnableToGenerateToken)
		return
	}
	// // w := http.ResponseWriter
	// data, _ := json.Marshal(fmt.Sprintf("authToken:%s,  userGroupId:%s", us.Token, us.UserGroupID))

	// // a := []byte(`{monster:[{basic:0,fun:11,count:262}],m:18}`)
	// fmt.Println(data)
	handleSuccessResponse(ctx, http.StatusOK, loginResult{AuthToken: us.Token, GroupId: us.UserGroupID, NameSpaceId: us.NamespaceID})
}

// Change password  godoc
// @Summary change password
// @Description Save  User
// @Tags User
// @Accept  json
// @Produce  json
// @Param  user body changePasswordParam true "change password changePasswordParam"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
//@Router /api/auth/changepassword [post]
func (ctrl *userController) ChangePassword(ctx *gin.Context) {

	var cpp changePasswordParam
	h := auth.Header{}

	if err := ctx.ShouldBindHeader(&h); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json header|ChangePassword") //, rest_errors.InvalidUserIdOrPassword)
		handleFaultResponse(ctx, saveErr, rest_errors.InvalidUserIdOrPassword)
		return
	}

	if err := ctx.ShouldBindJSON(&cpp); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json body|ChangePassword") //, rest_errors.InvalidUserIdOrPassword)
		// saveErr.SetErrorCode(rest_errors.InvalidUserIdOrPassword)
		// ctx.JSON(saveErr.Status(), saveErr)
		handleFaultResponse(ctx, saveErr, rest_errors.InvalidUserIdOrPassword)
		return
	}
	authToken, err := services.UserService.ChangePassword(cpp.EmailOrMobile, cpp.OldPassword, cpp.NewPassword, h.HeaderNS)
	if err != nil {
		//err.SetErrorCode(rest_errors.UnableToGenerateToken)
		handleFaultResponse(ctx, err, rest_errors.UnableToGenerateToken)
		return
	}
	handleSuccessResponse(ctx, http.StatusOK, authToken)
}

// Search godoc
// @Summary Get all models.User
// @Description get top 100 records models.User
// @Tags User
// @Produce  json
// @Param  start query int true "starting row"
// @Param  limit query int true "no of row limit"
// @Success 200 {array} models.User
// @Failure 404 {object} string
// @Failure 500 {object} string
//@Router /api/users/search [get]
//getting all  models.User
func (ctrl *userController) Search(ctx *gin.Context) {
	start := ctx.Query("start")
	limit := ctx.Query("limit")

	h := auth.Header{}

	if err := ctx.ShouldBindHeader(&h); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json header|Search") //, rest_errors.InvalidRequestHeader)
		handleFaultResponse(ctx, saveErr, rest_errors.InvalidRequestHeader)
		return
	}
	q := fmt.Sprintf("select * from %s where namespaceId:%s limit %s,%s", abstractimpl.UserTable, h.HeaderNS, start, limit)
	results, err := services.UserService.Search(q, h)
	if err != nil {
		saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("body json|%v", err)) //, rest_errors.InternalServerErrContactAdmin)
		logger.Error("Failed while search Users", saveErr)
		handleFaultResponse(ctx, err, rest_errors.InternalServerErrContactAdmin)
		return
	}
	handleSuccessResponse(ctx, http.StatusOK, results)
}

// Get godoc
// @Summary Get all models.User
// @Description Get models.User by id
// @Tags User
// @Produce  json
// @Param  id path string true "User ID"
// @Success 200 {array} models.User
// @Failure 404 {object} string
// @Failure 500 {object} string
//@Router /api/user/{id} [get]
//getting a record  models.User
func (ctrl *userController) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	h := auth.Header{}

	if err := ctx.ShouldBindHeader(&h); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json header|Search") //, rest_errors.InvalidRequestHeader)
		handleFaultResponse(ctx, saveErr, rest_errors.InvalidRequestHeader)
		return
	}

	us, err := services.UserService.Get(id, h)
	if err != nil {
		//	ctx.JSON(err.Status(), err)
		handleFaultResponse(ctx, err, rest_errors.InternalServerErrContactAdmin)
		return
	}
	handleSuccessResponse(ctx, http.StatusOK, us)
}

// Logout godoc
// @Summary Logout User
// @Description Logging out the specified user
// @Tags User
// @Produce  json
// @Success 200 {array} models.User
// @Failure 404 {object} string
// @Failure 500 {object} string
//@Router /api/auth/logout [get]
func (ctrl *userController) Logout(ctx *gin.Context) {
	h := auth.Header{}
	if err := ctx.ShouldBindHeader(&h); err != nil {
		saveErr := rest_errors.NewBadRequestError("invalid json header|Logout") //, rest_errors.InvalidUserIdOrPassword)
		handleFaultResponse(ctx, saveErr, rest_errors.InvalidUserIdOrPassword)
		return
	}

	if err := services.UserService.Logout(h.HeaderXauth); err != nil {
		//TODO: Add the appropriate error code
		//err.SetErrorCode(rest_errors.InternalServerErrContactAdmin)
		handleFaultResponse(ctx, err, rest_errors.InvalidAuthToken)
		return
	}
	handleSuccessResponse(ctx, http.StatusOK, h.HeaderXauth)
}

// Delete godoc
// @Summary Delete User Data
// @Description Delete User data matched by id
// @Tags User
// @Produce  json
// @Param  id path string true "User ID"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/user/{id} [delete]
// Deleting a record models.User
func (ctrl *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	h := auth.Header{}
	if err := ctx.ShouldBindHeader(&h); err != nil {
		deleteErr := rest_errors.NewBadRequestError(fmt.Sprintf("header|User|delete|%v", err)) //, rest_errors.InvalidUserIdOrPassword)
		logger.Error("Failed while delete", deleteErr)
		handleFaultResponse(ctx, deleteErr, rest_errors.InvalidUserIdOrPassword)
		return
	}
	err := services.UserService.Delete(id, h)
	if err != nil {
		//err.SetErrorCode(rest_errors.InternalServerErrContactAdmin)
		logger.Error("Failed while deleting User", err)
		handleFaultResponse(ctx, err, rest_errors.InternalServerErrContactAdmin)
		return
	}
	handleSuccessResponse(ctx, http.StatusOK, "Succesfully Deleted")
}
