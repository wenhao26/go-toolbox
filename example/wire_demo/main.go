package main

import (
	"fmt"
)

type Message struct {
	msg string
}

type Greeter struct {
	Message Message
}

type Event struct {
	Greeter Greeter
}

func NewMessage(msg string) Message {
	return Message{msg: msg}
}

func NewGreeter(message Message) Greeter {
	return Greeter{Message: message}
}

func NewEvent(greeter Greeter) Event {
	return Event{Greeter: greeter}
}

func (g Greeter) Greet() Message {
	return g.Message
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	//message := NewMessage("wire demo")
	//greeter := NewGreeter(message)
	//event := NewEvent(greeter)
	event := InitEvent("wire demo")

	event.Start()
}
