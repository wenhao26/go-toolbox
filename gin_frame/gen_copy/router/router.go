package router

import (
	"github.com/gin-gonic/gin"

	"toolbox/gin_frame/gen_copy/controllers"
	"toolbox/gin_frame/gen_copy/middleware"
)

func RegisterRouters(engine *gin.Engine) {
	engine.GET("/", controllers.BaseController.Index)

	engine.Use(middleware.Request())
	r := engine.Group("/api")
	{
		v1 := r.Group("/v1")
		{
			trace := controllers.TraceController

			// 追踪访问行为
			v1.Group("/trace").GET("/access-behavior", trace.AccessBehavior)

			// todo...
		}
	}

}
