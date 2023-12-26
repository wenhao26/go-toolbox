package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":9501")
	if err != nil {
		fmt.Println("监听失败：", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("认可失败：", err)
			continue
		}

		// 启动协程处理连接
		go func(c net.Conn) {
			defer c.Close()
			for {
				reader := bufio.NewReader(c)

				var buf [128]byte
				n, err := reader.Read(buf[:])
				if err != nil {
					fmt.Println("从客户端读取失败：", err)
					break
				}
				recvStr := string(buf[:n])
				fmt.Println("收到客户端发来的数据：", recvStr)
				//_, _ = c.Write([]byte(recvStr))
			}
		}(conn)

		go func(c net.Conn) {
			t := time.NewTicker(3e9)

			for {
				select {
				case <-t.C:
					_, _ = c.Write([]byte("[Sever]Hello TCP!"))
				default:

				}

			}
		}(conn)
	}

}
