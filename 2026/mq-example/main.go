package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 统一封装结构体
type RabbitMQ struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	QueueName string
}

// NewRabbitMQ 创建并初始化 RabbitMQ 实例，自动完成交换机和队列的声明与绑定
func NewRabbitMQ(url, exchange, queue, routingKey string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("create channel failed: %w", err)
	}

	// 声明直连交换机
	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare exchange failed: %w", err)
	}

	// 声明持久化队列
	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare queue failed: %w", err)
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("bind queue failed: %w", err)
	}

	return &RabbitMQ{
		Conn:      conn,
		Channel:   ch,
		QueueName: queue,
	}, nil
}

// Publish 生产者发送消息方法
func (mq *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return mq.Channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 消息持久化，防止重启丢失
			Body:         body,
		},
	)
}

// Consume 消费者启动方法，根据传入的数量拉起多个常驻协程并发消费
func (mq *RabbitMQ) Consume(consumerTag string, workerCount int, wg *sync.WaitGroup) (<-chan amqp.Delivery, error) {
	// 设置每个协程的预取上限，防止单个协程积压过多消息，起到全自动流控作用
	err := mq.Channel.Qos(workerCount*2, 0, false)
	if err != nil {
		return nil, fmt.Errorf("set qos failed: %w", err)
	}

	// 注册消费者，获取原生数据通道
	deliveryCh, err := mq.Channel.Consume(
		mq.QueueName,
		consumerTag,
		false, // autoAck: false 开启手动确认，保证干完活再发送确认信号
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("start consume failed: %w", err)
	}

	log.Printf("成功开启消费监听，正在拉起 %d 个常驻协程...\n", workerCount)

	// 并发拉起指定数量的独立消费协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("消费协程 [%d] 已就位，进入常驻阻塞状态。\n", workerID)

			// 基于通道原生特性，无消息时协程自动挂起休眠，不占 CPU
			for msg := range deliveryCh {
				handleBusiness(workerID, msg.Body)

				// 业务处理成功后，手动回复 ACK 确认
				_ = msg.Ack(false)
			}

			log.Printf("消费协程 [%02d] 安全消化完在途存量，平滑退出。\n", workerID)
		}(i)
	}

	return deliveryCh, nil
}

// Close 关闭连接和信道，释放资源
func (mq *RabbitMQ) Close() {
	if mq.Channel != nil {
		_ = mq.Channel.Close()
	}
	if mq.Conn != nil {
		_ = mq.Conn.Close()
	}
	log.Println("RabbitMQ 物理连接与信道已安全释放。")
}

// handleBusiness 核心业务处理函数，接收消费到的原始字节数据进行处理
func handleBusiness(workerID int, data []byte) {
	fmt.Printf(" -> [Worker %02d] 正在解析并存储高频行情: %s\n", workerID, string(data))
	time.Sleep(100 * time.Millisecond) // 模拟业务落库或 API 耗时
}

// startMockProducer 独立的生产者协程，高频生成测试行情数据，随 ctx 取消而退出
func startMockProducer(ctx context.Context, mq *RabbitMQ, exchange, routingKey string) {
	ticker := time.NewTicker(10 * time.Millisecond) // 每 20 毫秒推送一条
	defer ticker.Stop()

	var sequence int64 = 0
	log.Println("上游独立模拟生产者协程已成功拉起。")

	for {
		select {
		case <-ctx.Done():
			log.Println("接收到关闭通知，模拟生产者协程安全断流退出。")
			return
		case <-ticker.C:
			sequence++
			payload := fmt.Sprintf(`{"symbol":"BTCUSDT","price":68320.75,"seq":%d}`, sequence)
			err := mq.Publish(ctx, exchange, routingKey, []byte(payload))
			if err != nil {
				log.Printf("生产者消息投递失败: %v\n", err)
			}
		}
	}
}

func main() {
	log.Println("分布式常驻清洗进程开始初始化...")

	mqURL := "amqp://admin:admin@localhost:5672/"
	exchange := "mq_example_exchange"
	queue := "mq_example_queue"
	routingKey := "mq_example_routing_key"
	workerCount := 40
	consumerTag := "ticker_worker"

	// 全局上下文，用于通知异步生产者停止推流
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 初始化 MQ
	mq, err := NewRabbitMQ(mqURL, exchange, queue, routingKey)
	if err != nil {
		log.Fatalf("初始化 RabbitMQ 失败: %v", err)
	}
	defer mq.Close()

	// 启动常驻协程独立消费
	var wg sync.WaitGroup
	_, err = mq.Consume(consumerTag, workerCount, &wg)
	if err != nil {
		log.Fatalf("启动消费引擎失败: %v", err)
	}

	// 启动独立协程模拟生产数据
	go startMockProducer(ctx, mq, exchange, routingKey)

	// 挂载物理信号守卫，阻塞主线程，保持常驻进程后台运行
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	caughtSignal := <-signalChan
	log.Printf("监听到物理退出信号 [%v]，触发优雅停机...\n", caughtSignal)

	// 先让独立的生产者协程闭嘴断流
	cancel()

	// 向 MQ 宣告取消消费者监听，停止接收新消息，并死等 30 个常驻协程把手头的活干完
	log.Println("正在向 RabbitMQ 发送 Cancel 断流宣告，并等待在途协程收尾...")
	_ = mq.Channel.Cancel(consumerTag, false)
	wg.Wait()

	log.Println("进程平滑下电成功。")
}
