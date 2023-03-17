package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8899,
	})
	if err != nil {
		fmt.Println("net.ListenUDP error : ", err)
		return
	}
	defer listen.Close()

	for {
		var data [1024]byte

		// 接收数据报文
		n, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println("listen.ReadFromUDP error:", err)
			continue
		}
		fmt.Printf("data == %v  , addr == %v , count == %v\n", string(data[:n]), addr, n)

		// 将数据发给客户端
		_, err = listen.WriteToUDP(data[:n], addr)
		if err != nil {
			fmt.Println("listen.WriteToUDP error:", err)
			continue
		}
	}

}
