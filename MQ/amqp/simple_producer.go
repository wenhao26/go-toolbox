package MQ

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	exchange := "test.20220820"
	//queue := "test.queue_01"
	routingKey := "test.20220820_key"

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", "admin", "admin", "192.168.1.202", 5672)
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(fmt.Errorf("MQ构建连接失败: %s \n", err))
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Errorf("MQ频道创建失败: %s \n", err))
	}

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("声明交换机失败: %s \n", err))
	}
	defer ch.Close()

	/*q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("未能声明队列: %s \n", err))
	}*/

	for {
		data := map[string]string{
			"_date": time.Now().String(),
		}
		message, _ := json.Marshal(data)
		err = ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
			//DeliveryMode: amqp.Persistent, // 消息持久化
			ContentType: "text/plain",
			Body:        message,
		})
		if err != nil {
			fmt.Println("推送消息异常：", err)
		}
		fmt.Printf(" [x] Sent %s\n", message)
		time.Sleep(time.Millisecond * 100)
	}

}
