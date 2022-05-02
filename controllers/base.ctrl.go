package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
)

const (
	okStatus = "ok"
)

func handleFaultResponse(ctx *gin.Context, err rest_errors.RestErr, errCode rest_errors.InternalErrors) {
	//TODO error message needs to be pushed ui
	ctx.JSON(err.Status(), gin.H{"errCode": errCode, "status": err.Status(), "errorDesc": err.Message(), "datetime": date_utils.GetNowSearchFormat()})
}

func handleSuccessResponse(ctx *gin.Context, status int, data interface{}) {
	//TODO error message needs to be pushed ui server
	//logger.Debug(fmt.Sprintf("Success Resp %s", data))
	ctx.JSON(status, data)
}
