package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"toolbox/2025/websocket-chat/internal"
)

var (
	// 前端页面模板（简单缓存）
	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	// WebSocket 升级器，允许跨域（生产环境请严格控制）
	upgrader = websocket.Upgrader{
		// 允许所有来源。生产环境中应做限制，例如检查 Origin。
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	// 启动 hub（内部管理所有连接与广播）
	hub := internal.NewHub()
	go hub.Run()

	// 提供前端页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "template render error", http.StatusInternalServerError)
			return
		}
	})

	// 提供静态文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// WebSocket 连接点
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 使用 gorilla 升级连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}

		// 创建 client 并注册到 hub
		c := internal.NewClient(hub, conn)
		hub.Register <- c

		// 启动读写协程
		go c.WritePump()
		go c.ReadPump()
	})

	// 示例：在控制台可通过 HTTP 请求触发系统广播
	// 例如 curl -X POST "http://localhost:8080/system?msg=hello"
	http.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			http.Error(w, "missing msg", http.StatusBadRequest)
			return
		}
		// 系统消息结构（内部 hub 支持）
		internalMsg := internal.Message{
			Type:    internal.MessageTypeSystem,
			From:    "SYSTEM",
			Content: msg,
		}
		hub.Broadcast <- internalMsg
		_, _ = fmt.Fprintln(w, "ok")
	})

	addr := ":8080"
	log.Println("server listen:", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
