package main

import (
	"fmt"
	"time"

	"toolbox/simple-queue/mq"
)

var (
	topic = "test_topic"
)

// 单个主题测试
func OnceTopic() {
	m := mq.NewClient()
	defer m.Close()
	m.SetConditions(10)
	ch, err := m.Subscribe(topic)
	if err != nil {
		fmt.Println("subscribe failed")
		return
	}

	go OncePub(m)
	OnceSub(ch, m)
}

// 模拟定时推送
func OncePub(c *mq.Client) {
	t := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-t.C:
			currentDate := time.Now().String() + ":模拟推送了一条消息~"
			err := c.Publish(topic, currentDate)
			if err != nil {
				fmt.Println("publish message failed")
			}
		default:

		}
	}
}

// 接受订阅消息
func OnceSub(msg <-chan interface{}, c *mq.Client) {
	for {
		val := c.GetPayload(msg)
		fmt.Printf("get message is %s\n", val)
	}
}

// 多个主题测试
func ManyTopic() {
	m := mq.NewClient()
	defer m.Close()
	m.SetConditions(10)
	topic := ""
	for i := 0; i < 10; i++ {
		topic = fmt.Sprintf("test_topic_%d", i)
		go ManySub(m, topic)
	}
	ManyPub(m)
}

// 多个主题推送
func ManyPub(c *mq.Client) {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			for i := 0; i < 10; i++ {
				// 多个topic推送不同的消息
				topic := fmt.Sprintf("test_topic_%d", i)
				payload := fmt.Sprintf("多个主题推动_%d", i)
				err := c.Publish(topic, payload)
				if err != nil {
					fmt.Println("publish message failed")
				}
			}
		default:

		}
	}
}

func ManySub(c *mq.Client, topic string) {
	ch, err := c.Subscribe(topic)
	if err != nil {
		fmt.Printf("sub top:%s failed\n", topic)
	}

	for {
		val := c.GetPayload(ch)
		if val != nil {
			fmt.Printf("%s get message is %s\n", topic, val)
		}
	}
}

func main() {
	OnceTopic()
	//ManyTopic()
}
