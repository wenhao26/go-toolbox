package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// PriceUpdate 定义要发送的数据结构
type PriceUpdate struct {
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

// rootHandler 处理根路径，返回 HTML 文件
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// sseHandler 处理 SSE 连接
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// 设置 SSE 所需的 HTTP 头部
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取 ResponseWriter 的 Flusher 接口，用于强制发送缓冲区数据
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// 初始化一个虚拟价格
	currentPrice := 60000.00

	// 使用 select 结构，监听客户端连接关闭信号（r.Context().Done()）
	// 并启动一个无限循环来推送数据
	for {
		select {
		case <-r.Context().Done():
			// 客户端连接断开或浏览器关闭，退出 Goroutine
			fmt.Println("Client disconnected")
			return

		default:
			// 模拟价格波动
			// 价格在 +/- 100 范围内随机波动
			currentPrice += (rand.Float64() - 0.5) * 200.00

			// 构建数据
			update := PriceUpdate{
				Price:     currentPrice,
				Timestamp: time.Now().Format("15:04:05"),
			}

			// 序列化为 JSON 字符串
			jsonData, _ := json.Marshal(update)

			// 格式化并写入 SSE 响应
			// SSE 格式: "data: <JSON数据>\n\n"
			_, _ = fmt.Fprintf(w, "data: %s\n\n", jsonData)

			// 强制发送缓冲区数据到客户端
			flusher.Flush()

			// 暂停 1 秒，控制推送频率
			time.Sleep(1 * time.Second)
		}
	}

}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/events", sseHandler)

	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)

	// 启动 HTTP 服务器
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
