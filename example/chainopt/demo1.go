package main

import (
	"fmt"
)

type Handler struct {
	props interface{}

	err error
}

func (h *Handler) Step1() *Handler {
	if h.err != nil {
		return h
	}
	return h
}

func (h *Handler) Step2() *Handler {
	if h.err != nil {
		return h
	}
	return h
}

func (h *Handler) Err() error {
	return h.err
}

func main() {
	h := &Handler{}
	if err := h.Step1().Step2().Err(); err != nil {
		panic(err)
	}
	fmt.Println("OK")
}
