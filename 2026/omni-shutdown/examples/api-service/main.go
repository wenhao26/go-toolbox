package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	"toolbox/2026/omni-shutdown/pkg/shutdown"
)

type MockDatabase struct{}

func (m *MockDatabase) Close(ctx context.Context) error {
	log.Println("[DB-Infrastructure] 正在有序向 MySQL/Redis 连接池发出注销宣告...")
	time.Sleep(800 * time.Millisecond)
	log.Println("[DB-Infrastructure] 数据库网络 Socket 连接池已平滑安全释放。")
	return nil
}

func userAPIHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(200 * time.Millisecond) // 模拟轻量业务处理
	_, err := w.Write([]byte(`{"code":200,"msg":"success"}`))
	if err != nil {
		return
	}
}

func main() {
	// 场景一：经典高频 API 服务（限时超时强杀模式）
	// 适用于对外提供 Web、RPC 接口的微服务场景。
	// 其特征是不允许进程无限期挂起，要求在收到停机指令的限定时间内（如 5 秒）断流并快速退场。
	log.Println("[API-App] 正在初始化企业级高频 API 关卡服务...")

	// 初始化引擎，设定 5 秒的绝对硬超时时间
	engine := shutdown.New(shutdown.WithTimeout(5 * time.Second))

	// 模拟注册下游基础设施资源（先注册，最后关闭）
	db := &MockDatabase{}
	engine.RegisterShutdown(db.Close)

	// 构建标准 HTTP 服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(userAPIHandler),
	}

	// 注册上游流量断流钩子（后注册，率先关闭）
	engine.RegisterShutdown(server.Shutdown)

	// 异步拉起常驻监听
	go func() {
		log.Println("[API-App] HTTP 路由已成功绑定端口 :8080")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[API-App] 服务由于非正常原因发生崩溃: %v", err)
		}
	}()

	// 主协程移交控制权，阻塞守护
	code := engine.WaitListen()
	os.Exit(code)
}
