package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func process(conn net.Conn) {
	// 关闭连接
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [256]byte

		// 读取数据
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("reader.Read error:", err)
			break
		}

		recvData := time.Now().String() + ":" + string(buf[:n])
		fmt.Println("receive data", recvData)
		// 将数据再发给客户端
		_, _ = conn.Write([]byte(recvData))
	}
}

func main() {
	// 监听tcp
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}

	for {
		// 建立连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept error:", err)
			continue
		}
		fmt.Println("client connection:", conn.LocalAddr())

		// 开启协程去处理连接
		go process(conn)
	}
}
