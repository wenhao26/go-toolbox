package main

import (
	"log"
	"toolbox/2026/cron-jobs/internal/handlers"
	"toolbox/2026/cron-jobs/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化核心管理器
	cronMgr := service.NewCronManager()

	// 初始化Gin
	r := gin.Default()

	// 注册路由
	h := &handlers.TaskHandler{Mgr: cronMgr}
	h.RegisterRoutes(r)

	// 启动服务
	log.Println("中台任务管理系统启动在 :8088")
	if err := r.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}
