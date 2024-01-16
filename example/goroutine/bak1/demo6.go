package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

type Service struct {
	capacity     int
	tasks        chan *TaskData
	numThread    int
	closeChans   chan struct{}
	stopFlag     int32
	loopStopChan chan struct{}
}

type TaskData struct {
}

func NewService(capacity int) *Service {
	service := &Service{}

	service.capacity = capacity
	service.numThread = runtime.NumCPU() * 2
	service.tasks = make(chan *TaskData, capacity)
	service.stopFlag = 0
	service.closeChans = make(chan struct{}, service.numThread)
	service.loopStopChan = make(chan struct{})

	return service
}

func (this *Service) Stop() {
	atomic.StoreInt32(&this.stopFlag, 1)
	<-this.loopStopChan
	close(this.tasks)
	for i := 0; i < this.numThread; i++ {
		<-this.closeChans
	}
}

func (this *Service) Run() {
	for i := 0; i < this.numThread; i++ {
		go this.run(i)
	}
	go this.LoopConsume()
}

func (this *Service) run(i int) {
	fmt.Println("go run:", i)
loop:
	for {
		select {
		case task, ok := <-this.tasks:
			if ok {
				//#TODO process
				fmt.Println("process", task)
			} else {
				break loop
			}
		}
	}
	this.closeChans <- struct{}{}
}

func (this *Service) LoopConsume() {
	fmt.Println("loop")
	for atomic.LoadInt32(&this.stopFlag) == 0 {
		//TODO ReadData
		task := &TaskData{}
		this.tasks <- task

		fmt.Println("consume.")
		time.Sleep(time.Second * 2)
	}
	this.loopStopChan <- struct{}{}
}

func main() {
	service := NewService(100)
	go service.Run() //启动程序处理

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-c //等待关闭信号
	fmt.Println(s)
	service.Stop() //关闭service
	fmt.Println("exit :D")
}
