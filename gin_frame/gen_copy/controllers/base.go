package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Success = 1
	Failed  = 0
)

type Controller struct {
}

var BaseController = &Controller{}

func (*Controller) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Index")
}

func (*Controller) Success(ctx *gin.Context, msg string, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":  Success,
		"msg":   msg,
		"data":  data,
		"_time": time.Since(ctx.MustGet("_t").(time.Time)).String(),
	})
}

func (*Controller) Failed(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
