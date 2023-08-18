package mq

import (
	"errors"
	"sync"
	"time"
)

// 定义相关接口
type Broker interface {
	// 进行消息的推送，有两个参数即topic、msg，分别是订阅的主题、要传递的消息
	publish(topic string, msg interface{}) error

	// 消息的订阅，传入订阅的主题，即可完成订阅，并返回对应的channel通道用来接收数据
	subscribe(topic string) (<-chan interface{}, error)

	// 取消订阅，传入订阅的主题和对应的通道
	unsubscribe(topic string, sub <-chan interface{}) error

	// 关闭消息队列
	close()

	// 这个属于内部方法，作用是进行广播，对推送的消息进行广播，保证每一个订阅者都可以收到
	broadcast(msg interface{}, subscribers []chan interface{})

	// 这里是用来设置条件，条件就是消息队列的容量，这样我们就可以控制消息队列的大小了
	setConditions(capacity int)
}

type BrokerImpl struct {
	exit     chan bool
	capacity int

	topics       map[string][]chan interface{} // key:topic value:queue
	sync.RWMutex                               // 同步锁。读写锁，这里是为了防止并发情况下，数据的推送出现错误，所以采用加锁的方式进行保证
}

func NewBroke() *BrokerImpl {
	return &BrokerImpl{
		exit:   make(chan bool),
		topics: make(map[string][]chan interface{}),
	}
}

func (b *BrokerImpl) publish(topic string, msg interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:

	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	b.broadcast(msg, subscribers)
	return nil
}

func (b *BrokerImpl) subscribe(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker closed")
	default:

	}

	ch := make(chan interface{}, b.capacity)
	b.Lock()
	b.topics[topic] = append(b.topics[topic], ch)
	b.Unlock()
	return ch, nil
}

func (b *BrokerImpl) unsubscribe(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:

	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()

	if !ok {
		return nil
	}

	// 删除订阅者
	b.Lock()
	var newSubs []chan interface{}
	for _, subscriber := range subscribers {
		if subscriber == sub {
			continue
		}
		newSubs = append(newSubs, subscriber)
	}

	b.topics[topic] = newSubs
	b.Unlock()
	return nil
}

func (b *BrokerImpl) close() {
	select {
	case <-b.exit:
		return
	default:
		close(b.exit)
		b.Lock()
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}
}

func (b *BrokerImpl) broadcast(msg interface{}, subscribers []chan interface{}) {
	count := len(subscribers)
	concurrency := 1

	switch {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}

	pub := func(start int) {
		// 采用Timer，而不是使用time.After
		// 原因：time.After会产生内存泄漏，在计时器触发之前，垃圾回收不会回收Timer
		idleDuration := 5 * time.Millisecond
		idleTimeout := time.NewTimer(idleDuration)
		defer idleTimeout.Stop()

		for j := start; j < count; j += concurrency {
			if !idleTimeout.Stop() {
				select {
				case <-idleTimeout.C:
				default:

				}
			}
			idleTimeout.Reset(idleDuration)

			select {
			case subscribers[j] <- msg:
			case <-idleTimeout.C:
			case <-b.exit:
				return
			}
		}
	}

	for i := 0; i < concurrency; i++ {
		go pub(i)
	}
}

func (b *BrokerImpl) setConditions(capacity int) {
	b.capacity = capacity
}
