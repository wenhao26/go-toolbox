package bean

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type Bean struct {
	Conn *beanstalk.Conn
}

func New() *Bean {
	conn, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		log.Fatalln(err)
	}
	//defer conn.Close()
	return &Bean{conn}
}

func (b *Bean) CloseBean() {
	b.Conn.Close()
}

// 生产者
func (b *Bean) Produce(fName, tubeName string, body []byte, priority uint32, delay, ttr time.Duration) (uint64, error) {
	if fName == "" || tubeName == "" {
		return 0, errors.New("管道为空")
	}

	b.Conn.Tube.Name = tubeName
	b.Conn.TubeSet.Name[tubeName] = true

	log.Println(fName, " [Producer] tubeName:", tubeName, " c.Tube.Name:", b.Conn.Tube.Name)

	// 参数分别为任务，优先级，延时时间，处理任务时间
	id, err := b.Conn.Put(body, priority, delay, ttr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("PUB#JOB-ID=", id)
	return 0, err
}

// 消费者
/*func (b *Bean) Consume1(tubeName string) {
	b.Conn.Tube.Name = tubeName
	b.Conn.TubeSet.Name[tubeName] = true

	for {
		// 为了防止某个consume长时间占用，设置了timeout
		id, body, err := b.Conn.Reserve(5 * time.Second)
		if err != nil {
			log.Println(err)
		}
		if id > 0 {
			log.Println("SUB#JOB-ID=", id, ";BODY=", string(body))
			_ = b.Conn.Delete(id)
		}
	}
}*/
func (b *Bean) Consume(fName, tubeName string) {
	b.Conn.Tube.Name = tubeName
	b.Conn.TubeSet.Name[tubeName] = true

	log.Println(fName, " [Consumer] tubeName:", tubeName, " c.Tube.Name:", b.Conn.Tube.Name)

	substr := "timeout"
	for {
		// 为了防止某个consume长时间占用，设置了timeout
		id, body, err := b.Conn.Reserve(3 * time.Second)
		if err != nil {
			if strings.Contains(err.Error(), substr) {
				log.Println(fName, " [Consumer] [", b.Conn.Tube.Name, "] err:", err, " id:", id)
			}
			continue
		}

		log.Println("SUB#JOB-ID=", id, ";BODY=", string(body))

		// 从队列中清掉
		err = b.Conn.Delete(id)
		if err != nil {
			log.Println(fName, " [Consumer] [", b.Conn.Tube.Name, "] Delete err:", err, " id:", id)
		} else {
			log.Println(fName, " [Consumer] [", b.Conn.Tube.Name, "] Successfully deleted. id:", id)
		}
	}
}

// 查看管道
func (b *Bean) WatchTubes() {
	arr, err := b.Conn.ListTubes()
	if err != nil {
		log.Println("[Example] err:", err)
	} else {
		for i, v := range arr {
			log.Println("[Example] ListTubes  i:", i, " v:", v)
			b.Conn.Tube.Name = v
			id, body, err := b.Conn.Reserve(5 * time.Second)
			if err != nil {
				log.Println("[Example] err:", err, " name:", b.Conn.Tube.Name)
				continue
			} else {
				log.Println("[Example] job:", id)
				log.Println("[Example] body:", string(body))
			}
		}
	}
}

// 管道状态
func (b *Bean) TubeStat() {
	log.Println("Tube(", b.Conn.Tube.Name, ") Stats:")

	m, err := b.Conn.Tube.Stats()
	if err != nil {
		log.Println("[tubeStatus] err:", err)
	} else {
		for k, v := range m {
			log.Println(k, " : ", v)
		}
	}
}
