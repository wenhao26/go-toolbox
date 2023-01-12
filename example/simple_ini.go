package main

import (
	"fmt"

	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("../conf/ini/my.ini")
	if err != nil {
		panic(err)
	}

	sec := cfg.Section("cent")
	addr := sec.Key("addr").String()
	fmt.Println(addr)
}
