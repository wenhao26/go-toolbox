package main

import (
	"fmt"
	"log"
	"toolbox/2026/gocron/internal/handler"
	"toolbox/2026/gocron/internal/service"

	"github.com/gin-gonic/gin"
)

// -- 入口主程序 --

func main() {
	fmt.Println("==================================================")
	fmt.Println("        🚀 Go-Cron 分布式秒级任务系统 启动中...      ")
	fmt.Println("==================================================")

	// 初始化 Redis
	redisAddr := "192.168.33.10:6379"
	service.InitRedis(redisAddr)
	fmt.Printf("[1/4] 🔗 正在连接 Redis: %s ... ✅\n", redisAddr)

	// 初始化调度器
	service.InitScheduler()
	fmt.Println("[2/4] ⚙️ 调度引擎已就绪 (支持秒级解析) ... ✅")

	// 加载Redis任务
	fmt.Println("[3/4] 📥 正在从持久化层同步历史任务...")
	tasks, err := service.LoadAllTasks()
	if err != nil {
		log.Fatalf("❌ 无法从 Redis 加载任务: %v", err)
	}

	activeCount := 0
	for _, task := range tasks {
		if task.Status == 1 {
			_ = service.GlobalTaskScheduler.AddOrUpdate(task)
			fmt.Printf("	-> 🟢 激活任务: [%s] Cron: [%s]\n", task.ID, task.Expr)
			activeCount++
		} else {
			fmt.Printf("	-> ⚪ 忽略任务: [%s] (已暂停)\n", task.ID)
		}
	}
	fmt.Printf("✨ 同步完成，共加载 %d 个任务，激活 %d 个。\n", len(tasks), activeCount)

	// 设置路由并启动
	fmt.Println("[4/4] 🌐 正在注册 HTTP 接口与静态页面...")

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.StaticFile("/", "./static/index.html")
	r.POST("/api/save", handler.SaveTask)
	r.POST("/api/delete", handler.DeleteTask)
	r.GET("/api/tasks", handler.GetTaskList)
	r.GET("/api/logs", handler.GetLogs)

	fmt.Println("==================================================")
	fmt.Printf(" ✅ 服务启动成功! 请访问: http://localhost:8989\n")
	fmt.Println("==================================================")

	if err := r.Run(":8989"); err != nil {
		log.Fatalf("❌ 服务启动失败: %v", err)
	}
}
