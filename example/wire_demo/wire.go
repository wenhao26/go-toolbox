package main

import (
	"github.com/google/wire"
)

func InitEvent(msg string) Event {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}
}
