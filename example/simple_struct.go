package main

import (
	"fmt"
)

type Options struct {
	Addr      string
	Port      int
	StorageDB string
}

func (opt *Options) New() {
	fmt.Println(opt.Addr)
}

func main() {
	opt := &Options{
		Addr:      "127.0.0.1",
		Port:      3306,
		StorageDB: "aws-s3",
	}
	opt.New()
}
