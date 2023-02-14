package controllers

import (
	"github.com/gin-gonic/gin"

	"toolbox/gin_frame/gen_copy/zlog"
)

type trace struct {
	*Controller
}

var TraceController = &trace{
	Controller: BaseController,
}

func (t *trace) AccessBehavior(ctx *gin.Context) {
	zlog.Info("Request-API:/trace/access-behavior")
	t.Success(ctx, "响应成功", nil)
}
