package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9502,
	})
	if err != nil {
		fmt.Println("监听失败：", err)
		return
	}

	// 关闭连接
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭失败：", err)
		}
	}(listen)

	for {
		var data [1024]byte

		// 读取UDP数据
		count, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println("读取失败：", err)
			continue
		}
		fmt.Printf("data: %s add: %v count: %v \n\n", string(data[0:count]), addr, count)

		// 发送数据
		_, err = listen.WriteToUDP([]byte("[server]hello udp"), addr)
		if err != nil {
			fmt.Println("发送失败：", err)
			continue
		}
	}

}
