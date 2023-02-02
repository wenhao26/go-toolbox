package main

import (
	"fmt"
)

type DoTask interface {
	task()
}

type Pub1 struct {
}

type Pub2 struct {
}

func (p Pub1) task() {
	fmt.Println("PUB1")
}

func main() {
	var p1 Pub1

	p1.task()
}
