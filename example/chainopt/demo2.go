package main

import (
	"fmt"
)

type Server struct {
	host string
	port string
}

type Option func(s *Server)

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{host: "127.0.0.1", port: "9521"}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func main() {
	//s := NewServer()
	s := NewServer(WithHost("localhost"), WithPort("88867"))
	fmt.Println(s)
}
