package main

import (
	"fmt"

	"github.com/mholt/archiver/v3"
)

// 解压
func main() {
	filename := "F:\\zip-test\\example.zip"
	err := archiver.Unarchive(filename, "F:\\zip-test\\example")
	if err != nil {
		panic(err)
	}
	fmt.Println("OK")
}
