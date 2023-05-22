package main

import (
	"fmt"

	"toolbox/utils"
)

func main() {
	email := "test@gmail.com"
	ret := utils.ValidateEmail(email)

	fmt.Println(ret)
}
