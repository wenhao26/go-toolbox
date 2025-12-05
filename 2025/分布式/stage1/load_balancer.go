// 实现一个轮询或随机的负载均衡器，将传入请求转发给多个后端 HTTP 服务实例
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Backend 结构体：表示一个后端服务实例
type Backend struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy // Go 标准库中的反向代理对象
}

// ==============================
// 2. 负载均衡器核心结构 (Round Robin)
// ==============================

// LoadBalancer 结构体：管理后端列表和轮询状态
type LoadBalancer struct {
	Backends []Backend
	Current  int        // 当前轮询到的索引
	Mutex    sync.Mutex // 保护 Current 字段，避免并发访问问题
}

// NextBackend 实现轮询策略 (Round Robin)
func (lb *LoadBalancer) NextBackend() *Backend {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	// 轮询逻辑：Current % 列表长度
	// 每次调用，Current 都会递增，然后对 Backends 列表的长度取模，实现循环访问。
	index := lb.Current % len(lb.Backends)
	lb.Current = (lb.Current + 1) % len(lb.Backends) // 更新下一个索引

	return &lb.Backends[index]
}

// ServeHTTP 是 http.Handler 接口的实现，负责处理传入请求
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. 使用轮询策略选择下一个后端
	backend := lb.NextBackend()

	log.Printf("[代理请求] 转发请求到后端: %s", backend.URL)

	// 2. 将请求转发给选定的后端
	backend.Proxy.ServeHTTP(w, r)
}

// startBackendServer 启动 HTTP 服务
func startBackendServer(port string, id int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[服务 %d] 接收到请求", id)
		time.Sleep(200 * time.Millisecond)
		_, _ = fmt.Fprintf(w, "服务 %d，运行在端口 %s", id, port)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("启动后端服务 %d 在 http://localhost:%s\n", id, port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("后端服务 %d 启动失败: %v", id, err)
	}
}

func main() {
	backendPorts := []string{"9001", "9002", "9003"} // 后端服务端口列表

	// 启动所有后端服务，并发启动
	for i, port := range backendPorts {
		go startBackendServer(port, i+1)
	}

	// 等待服务启动 (实际项目中应使用更健壮的健康检查)
	time.Sleep(time.Second)

	// 初始化 Load Balancer
	var backends []Backend
	for _, port := range backendPorts {
		backendURL, _ := url.Parse("http://localhost:" + port)

		// 创建反向代理，并设置目标 URL
		proxy := httputil.NewSingleHostReverseProxy(backendURL)

		backends = append(backends, Backend{
			URL:   backendURL,
			Proxy: proxy,
		})
	}

	lb := &LoadBalancer{
		Backends: backends,
		Current:  0,
	}

	// 启动负载均衡器 (代理服务)
	lbPort := "8080"
	fmt.Printf("\n===== 负载均衡器启动在 http://localhost:%s =====\n", lbPort)
	log.Fatal(http.ListenAndServe(":"+lbPort, lb))
}
