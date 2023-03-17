package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("client start...")

	for i := 0; i < 30; i++ {
		msg := `Hello TCP!`
		conn.Write([]byte(msg))
	}

	fmt.Println("数据发送完了~")
}
