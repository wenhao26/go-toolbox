package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func process2(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var buf [2048]byte

	for {
		n, err := reader.Read(buf[:])
		// 如果客户端关闭，则退出本协程
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("reader.Read error:", err)
			break
		}

		recvStr := string(buf[:n])
		// 打印收到的数据
		fmt.Printf("接收到的数据：%s\n\r", recvStr)
	}
}

// 出现黏包问题
func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listen.Close()

	fmt.Println("server start...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept error:", err)
			continue
		}

		go process2(conn)
	}
}
