package MQ

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func main() {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", "admin", "admin", "192.168.1.202", 5672)
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(fmt.Errorf("MQ构建连接失败: %s \n", err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Errorf("MQ频道创建失败: %s \n", err))
	}

	go func() {
		messages, err := ch.Consume("test.queue_01", "TEST", true, false, false, false, nil)
		if err != nil {
			panic(fmt.Errorf("消费消息异常: %s \n", err))
		}

		for message := range messages {
			log.Printf("[MESSAGE]=%s", string(message.Body))
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
