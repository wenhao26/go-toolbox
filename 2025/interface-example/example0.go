package main

import (
	"fmt"
)

func stringType(v interface{}) {
	if s, ok := v.(string); ok {
		fmt.Println("string：", s)
	} else {
		fmt.Println("not a string")
	}
}

func main() {
	stringType("isOK")
	stringType(1688)
}
