package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/auth"
	"github.com/ramnkl16/ez-search/global"
)

var (
	FileController fileControllerInteface = &fileController{}
)

type fileControllerInteface interface {
	upload(ctx *gin.Context)
	// download(ctx *gin.Context)
	// delete(ctx *gin.Context)
	// getFileNames(ctx *gin.Context)
	RegisterRouter(rout *gin.Engine)
}

type fileController struct{}

func (ctrl *fileController) RegisterRouter(rout *gin.Engine) {
	rout.POST("/api/files/upload", ctrl.upload)
	// rout.PUT("/api/appLog", ctrl.Update)
	// rout.DELETE("/api/appLog/:id", ctrl.Delete)
	// rout.GET("/api/appLog/:id", ctrl.Get)
	// rout.GET("/api/appLogs/search", ctrl.Search)
	//rout.GET("/web-ui/appLogs/search", ctrl.Search)
}

type fileJson struct {
	UploadPath  string                `form:"path" binding:"required"`
	Description string                `form:"description"`
	fileData    *multipart.FileHeader `form:"file"`
}

func (ctrl *fileController) upload(c *gin.Context) {
	var testmultipart fileJson
	err := c.Bind(&testmultipart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "error": err.Error()})
	}
	//fmt.Println(testmultipart.UploadPath)
	//fmt.Println(testmultipart.Description)

	form, _ := c.MultipartForm()
	files := form.File["file"]
	//files := form.File["upload[]"]
	ns := auth.GetNamespace(c.Request)
	for _, file := range files {
		fmt.Println(file.Filename)
		dirName := filepath.Join(global.WorkingDir, ns, testmultipart.UploadPath)
		fullFileName := filepath.Join(global.WorkingDir, ns, testmultipart.UploadPath, file.Filename)
		//fmt.Println(dirName)
		//fmt.Println(fullFileName)
		_, err := os.Stat(dirName)
		if os.IsNotExist(err) {
			os.MkdirAll(dirName, os.ModePerm)
		}
		c.SaveUploadedFile(file, fullFileName)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusText(http.StatusOK), "data": "xx"})
}
