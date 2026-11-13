package main

import (
	"fmt"
)

// Server 定义一个服务器结构体
type Server struct {
	Host     string
	Port     int
	MaxToken int
}

// Option 定义一个修改 Server 配置的函数类型
type Option func(server *Server)

// WithPort 设置自定义端口的选项函数
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

// WithMaxToken 设置最大 Token 数的选项函数
func WithMaxToken(maxToken int) Option {
	return func(s *Server) {
		s.MaxToken = maxToken
	}
}

// NewServer 使用函数选项模式创建对象
func NewServer(host string, opts ...Option) *Server {
	server := &Server{
		Host:     host,
		Port:     8080,
		MaxToken: 1024,
	}

	// 依次应用用户传入的所有选项
	for _, opt := range opts {
		opt(server)
	}
	return server
}

// - 函数选项模式
// 场景：当你需要写一个结构体（比如一个 HTTP Server、RPC 客户端或自定义连接池），它有十几个配置项，但绝大多数都有合理的默认值。用户在创建时只想修改其中的 1~2 个。
//
// 传统痛点：如果用构造函数传参，参数会变得巨长无比；如果用结构体传参，又很难区分“用户没传”还是“用户传了零值”。
//
// 惊艳之处：Go 社区标配的函数选项模式，完美兼顾了向后兼容性、代码可读性和扩展性。
//
// 学习收获：这是 Go 语言中编写优雅 API、SDK 的高级设计模式。Docker、Etcd 等开源巨作的代码库中到处都是这种写法。

func main() {
	s1 := NewServer("127.0.0.1")
	fmt.Println(s1)

	s2 := NewServer("127.0.0.1", WithPort(8081), WithMaxToken(2048))
	fmt.Println(s2)
}
