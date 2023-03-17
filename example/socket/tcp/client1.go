package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 键入数据
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 读取用户输入
		input, _ := inputReader.ReadString('\n')
		// 截断
		inputInfo := strings.Trim(input, "\r\n")
		// 读取到用户输入q或者Q就退出
		if strings.ToUpper(inputInfo) == "Q" {
			return
		}

		// 将输入的数据发送给服务端
		_, err = conn.Write([]byte(inputInfo))
		if err != nil {
			return
		}

		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("conn.Read error:", err)
			return
		}
		fmt.Println("接收服务端消息：", string(buf[:n]))
	}
}
