package main

import (
	"fmt"

	"toolbox/AES"
)

const (
	KEY = "8sDKrFK2ojyqAf7fkXoHI6h7OcGTMqgw"
)

func main() {
	content := []byte("AES-CEB模式加/解密!")

	ceb := AES.New(KEY)
	ret := ceb.Encrypt(content)

	fmt.Println(ret)
	fmt.Println("---------------")

	ret1 := ceb.Decrypt(ret)
	fmt.Println(ret1)
}
