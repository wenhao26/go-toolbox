package main

import (
	"fmt"
)

// event 事件数据
type event struct {
	ID   string
	Name string
}

// mockUser 模拟用户数据
func mockUser(userCh chan<- string, done chan<- bool) {
	for i := 0; i < 10; i++ {
		userCh <- fmt.Sprintf("username_%d", i)
	}

	close(userCh)
	done <- true
}

// mockEvent 模拟事件数据
func mockEvent(eventCh chan<- event, done chan<- bool) {
	for i := 0; i < 20; i++ {
		eventCh <- event{
			ID:   fmt.Sprintf("id_%d", i),
			Name: fmt.Sprintf("name_%d", i),
		}
	}

	close(eventCh)
	done <- true
}

func main() {
	userChan := make(chan string)
	userDone := make(chan bool)
	eventChan := make(chan event)
	eventDone := make(chan bool)

	go mockUser(userChan, userDone)
	go mockEvent(eventChan, eventDone)

	count := 0

	for count < 2 {
		select {
		case val, ok := <-userChan:
			if !ok {
				fmt.Println("userChan已关闭")
			}
			fmt.Println(val)
		case val, ok := <-eventChan:
			if !ok {
				fmt.Println("eventChan已关闭")
			}
			fmt.Println(val)
		case <-userDone:
			count++
		case <-eventDone:
			count++
		}
	}

	fmt.Println("Done")
}
