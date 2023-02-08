package main

import (
	"fmt"

	rdb2 "toolbox/extend/rdb"
)

func main() {
	rdb := rdb2.NewRdb(&rdb2.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       2,
	})

	rdb.Set("test", "12345678", 0)

	fmt.Println("DONE")
}
