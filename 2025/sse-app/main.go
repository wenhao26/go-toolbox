package main

import (
	"encoding/json"
	"fmt"
	"math"
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

const (
	BasePrice     = 60000.00 // 基础价格（中心值）
	Amplitude     = 5000.00  // 波动幅度
	PeriodSeconds = 86400.0  // 波动周期：24小时 (86400秒)
)

// calculateDeterministicPrice 根据当前时间戳计算确定的模拟价格
func calculateDeterministicPrice() float64 {
	// 获取总的 Unix 时间戳（秒）
	t := float64(time.Now().Unix())

	// 将总时间戳限制在 0 到 PERIOD_SECONDS 之间
	// math.Mod(t, PERIOD_SECONDS) 得到今天是这个周期（24小时）的第几秒
	tMod := math.Mod(t, PeriodSeconds)

	// 将时间映射到三角函数周期 (0 到 2π)
	// phase 是当前时间戳在周期中的位置
	phase := (tMod / PeriodSeconds) * 2 * math.Pi

	// 使用余弦函数计算价格波动
	// 余弦函数在 [-1, 1] 之间波动，保证价格变化自然平滑
	oscillation := math.Cos(phase)

	// 价格 = 基础价格 + 波动幅度 * 波动值
	price := BasePrice + Amplitude*oscillation

	// 确保两位小数精度
	return math.Round(price*100) / 100
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
			currentPrice := calculateDeterministicPrice()

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
