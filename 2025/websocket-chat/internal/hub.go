package internal

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Hub 管理所有客户端连接，广播消息，注册/注销客户端。
// 设计思想：使用 channel 进行并发安全的事件驱动管理（类似 chat example）。

// MessageType 表示消息类型
type MessageType string

const (
	MessageTypeText   MessageType = "text"   // 普通用户消息
	MessageTypeSystem MessageType = "system" // 系统消息（由服务端发送）
	MessageTypeEvent  MessageType = "event"  // 事件消息，如登录/退出
)

// Message 是所有广播消息的统一结构
type Message struct {
	Type    MessageType `json:"type"`
	From    string      `json:"from"`    // 发送者，登录时为用户名，系统消息为 SYSTEM
	Content string      `json:"content"` // 文本内容
	Time    int64       `json:"time"`    // unix 时间戳（秒）
}

// Client 代表一个 websocket 连接
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message // 待发送到客户端的消息

	// 用户名。登录后由前端发送设置用户名事件（见 ReadPump）
	name string
}

// Hub 定义
type Hub struct {
	// 注册/注销 client
	Register   chan *Client
	Unregister chan *Client

	// 广播消息
	Broadcast chan Message

	// 当前活跃客户端集合（key 为 *Client）
	clients map[*Client]bool
}

// NewHub 创建 hub
func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		clients:    make(map[*Client]bool),
	}
}

// Run 循环处理注册、注销、广播事件
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.clients[c] = true
			log.Println("新连接注册")
			// 连接建立后，等待客户端发送设置用户名的消息；在这里不直接广播


		case c := <-h.Unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				// 如果该客户端有用户名，则广播退出事件
				if c.name != "" {
					msg := Message{
						Type:    MessageTypeEvent,
						From:    c.name,
						Content: c.name + " 已退出",
						Time:    time.Now().Unix(),
					}
					log.Printf("%s 退出成功\n", c.name)
					h.broadcastToAll(msg)
				}
			}


		case msg := <-h.Broadcast:
			// 打印日志，对于系统/事件消息也统一打印
			switch msg.Type {
			case MessageTypeEvent:
				log.Printf("事件: %s -> %s\n", msg.From, msg.Content)
			case MessageTypeSystem:
				log.Printf("系统消息: %s\n", msg.Content)
			default:
				log.Printf("消息来自 %s: %s\n", msg.From, msg.Content)
			}
			h.broadcastToAll(msg)
		}
	}
}

// broadcastToAll 将消息推入每个客户端的发送通道（若阻塞则跳过）
func (h *Hub) broadcastToAll(msg Message) {
	for c := range h.clients {
		select {
		case c.send <- msg:
		default:
			// 若客户端发送缓冲已满，说明该客户端可能不可用，关闭连接并注销
			close(c.send)
			delete(h.clients, c)
		}
	}
}

// NewClient 创建 client
func NewClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  h,
		conn: conn,
		send: make(chan Message, 256), // 缓冲区，防止短时间大量广播导致阻塞
	}
}

// ReadPump 读取客户端发来的消息，处理登录/文本等动作
func (c *Client) ReadPump() {
	defer func() {
		// 在退出时反注册
		c.hub.Unregister <- c
		_ = c.conn.Close()
	}()

	// 设置读取限制和心跳（可扩展）
	c.conn.SetReadLimit(512)

	for {
		var msg Message
		// 这里我们简单以 JSON 格式的 Message 作为通信协议
		if err := c.conn.ReadJSON(&msg); err != nil {
			// 读取错误通常意味着连接关闭
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v\n", err)
			}
			break
		}

		// 根据消息类型做处理
		switch msg.Type {
		case MessageTypeEvent:
			// 事件类型用于登录（例如 content = "login"，from = username）
			// 约定：前端发送登录事件时，Message.From 为用户名，Message.Content = "login"
			if msg.Content == "login" && msg.From != "" {
				// 设置客户端用户名
				c.name = msg.From
				// 打印登录并广播
				log.Printf("%s 登录成功\n", c.name)
				event := Message{
					Type:    MessageTypeEvent,
					From:    c.name,
					Content: c.name + " 已登录",
					Time:    time.Now().Unix(),
				}
				c.hub.Broadcast <- event
			}

		case MessageTypeText:
			// 普通聊天消息，直接填充 From/Content/Time 后广播
			if c.name == "" {
				// 未登录的客户端不允许发送消息，可以选择忽略或返回错误消息
				warn := Message{
					Type:    MessageTypeSystem,
					From:    "SYSTEM",
					Content: "请先登录再发送消息",
					Time:    time.Now().Unix(),
				}
				c.send <- warn
				continue
			}
			chat := Message{
				Type:    MessageTypeText,
				From:    c.name,
				Content: msg.Content,
				Time:    time.Now().Unix(),
			}
			c.hub.Broadcast <- chat

		}
	}
}

// WritePump 将 hub 推送过来的消息写到 websocket
func (c *Client) WritePump() {
	defer func() {
		_ = c.conn.Close()
	}()

	for msg := range c.send {
		// 以 JSON 发送
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}
