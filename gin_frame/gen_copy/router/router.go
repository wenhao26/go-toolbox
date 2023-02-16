package router

import (
	"time"

	limitMaxAllowed "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"

	"toolbox/gin_frame/gen_copy/controllers"
	"toolbox/gin_frame/gen_copy/middleware"
)

func RegisterRouters(engine *gin.Engine) {
	engine.GET("/", controllers.BaseController.Index)

	// 基础请求中间件、同时最大允许访问
	engine.Use(middleware.Request(), limitMaxAllowed.MaxAllowed(100))

	// 限制访问速率
	engine.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		// 按客户端ip限制速率
		return c.ClientIP()
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		// 限制10个qps/clientIp，并允许最多10个令牌的突发，限制器活动时间为1小时
		return rate.NewLimiter(rate.Every(100*time.Millisecond), 10), time.Hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429)
	}))

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
