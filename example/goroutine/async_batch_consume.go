package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

/*
消息产生时间
	1：一秒一个信息
	2：一毫秒一个信息
	3：随机0~5秒一个信息
*/
var periodType = 2

// 消息管道大小限制
var msgChanLimit = 20

// 协程任务自定义ID，用来区分查看对应的任务批次，每处理一次自增1
var handleID int

// 最多堆积多少未处理信息，进行一批次处理
var batchNum = 6

// 在未满消息时，最长多久进行一批次处理
var batchTime = 3

// 最大的协程数量
var maxGoroutines = 10

// 发布消息
func publishMsg(msgChan chan string) {
	msgID := 0

	// 循环模拟发布消息
	for {
		msgID++
		body := "msg-id:" + strconv.Itoa(msgID) + ";add-time:" + time.Now().String()
		msgChan <- body

		switch periodType {
		case 1:
			time.Sleep(1 * time.Second) // 秒
		case 2:
			time.Sleep(1 * time.Millisecond) // 毫秒
		case 3:
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond) // 随机
		}
	}
}

// 处理消息
func handleMsg(model int, handleID int, msgSet []string, guard chan struct{}) {
	for i, v := range msgSet {
		fmt.Println("handle-id:", handleID, ";model:", model, ">>> idx:", strconv.Itoa(i), ";body:", v)
		time.Sleep(1500 * time.Millisecond) // 模拟具体处理消息消耗的时间
	}
	<-guard // 释放一个
}

// 无阻塞去接受消息
func unBlockRead(ch chan string) (msg string, err error) {
	select {
	case msg = <-ch:
		return msg, nil
	case <-time.After(time.Microsecond):
		return "", errors.New("nil")
		//return "消息超时", nil
	}
}

func main() {
	guard := make(chan struct{}, maxGoroutines) // 守护协程数量限制
	msgChan := make(chan string, msgChanLimit)  // 接收消息的大小
	msgSet := make([]string, 0)                 // 临时存放接收的消息集合
	step := 0                                   // 秒级别计数器

	// 生产消息
	go publishMsg(msgChan)

	// 开始处理
	for {
		if msg, err := unBlockRead(msgChan); err == nil {
			msgSet = append(msgSet, msg)
			if len(msgSet) == batchNum { // 达到处理数量
				handleID++
				guard <- struct{}{}
				go handleMsg(1, handleID, msgSet, guard) // 处理当前的msgSet
				msgSet = nil                             // 重置
				step = 0
			}
		} else {
			if step > batchTime && len(msgSet) > 0 { // 超时并且不为空
				handleID++
				guard <- struct{}{}
				go handleMsg(2, handleID, msgSet, guard)
				msgSet = nil
				step = 0
			} else {
				step++
				time.Sleep(1 * time.Second)
			}
		}
	}

	// 挂起主进程防止退出
	for {
		runtime.GC()
	}

}
