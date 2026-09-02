package producer

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"toolbox/2026/async-engine/engine"
	"toolbox/2026/async-engine/model"
)

// Producer 模拟高并发 API 访问并产生日志任务，通过令牌桶控制整体产出速率。
type Producer struct {
	cfg     Config
	engine  *engine.Engine
	limiter *Limiter

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	submitTotal   int64
	submitSuccess int64
}

// Snapshot 保存某一时刻的生产者指标快照。
type Snapshot struct {
	SubmitTotal   int64
	SubmitSuccess int64
}

// NewProducer 创建生产者实例，传入引擎用于任务投递。
func NewProducer(cfg Config, eng *engine.Engine) *Producer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Producer{
		cfg:     cfg,
		engine:  eng,
		limiter: NewLimiter(),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start 启动投递协程与流量阶段调度器。
func (p *Producer) Start() {
	for i := 0; i < p.cfg.WorkerNum; i++ {
		p.wg.Add(1)
		go p.dispatch(i)
	}
	go p.schedule()
}

// Stop 停止产生新任务，并等待所有投递协程退出。
func (p *Producer) Stop() {
	p.cancel()
	p.wg.Wait()
	fmt.Println("[producer] all dispatchers exited")
}

func (p *Producer) dispatch(id int) {
	defer p.wg.Done()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		if !p.limiter.Allow() {
			// 速率限制，短暂退让避免 CPU 空转
			time.Sleep(100 * time.Microsecond)
			continue
		}

		task := p.generateTask(rnd)
		atomic.AddInt64(&p.submitTotal, 1)

		if p.engine.Submit(task) {
			atomic.AddInt64(&p.submitSuccess, 1)
		} else {
			// 队列背压，短暂退让降低竞争
			time.Sleep(time.Millisecond)
		}
	}
}

func (p *Producer) schedule() {
	if len(p.cfg.Phases) == 0 {
		return
	}

	for _, phase := range p.cfg.Phases {
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		p.limiter.SetRate(phase.QPS)
		fmt.Printf("[producer] phase started: qps=%d duration=%s\n", phase.QPS, phase.Duration)

		timer := time.NewTimer(phase.Duration)
		select {
		case <-p.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}

	fmt.Println("[producer] all phases completed, holding last rate")
}

func (p *Producer) generateTask(rnd *rand.Rand) *model.Task {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	paths := []string{
		"/api/v1/users",
		"/api/v1/orders",
		"/api/v1/products",
		"/api/v1/payments",
		"/health",
	}
	statusCodes := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 404, 500, 503}

	return &model.Task{
		TraceID:     fmt.Sprintf("trace-%d-%d", time.Now().UnixNano(), rnd.Int63()),
		Timestamp:   time.Now(),
		Method:      methods[rnd.Intn(len(methods))],
		Path:        paths[rnd.Intn(len(paths))],
		Status:      statusCodes[rnd.Intn(len(statusCodes))],
		Latency:     time.Duration(rnd.Intn(100)) * time.Millisecond,
		ClientIP:    fmt.Sprintf("192.168.%d.%d", rnd.Intn(256), rnd.Intn(256)),
		PayloadSize: rnd.Intn(1024 * 1024),
	}
}

// Snapshot 原子读取当前生产者指标。
func (p *Producer) Snapshot() Snapshot {
	return Snapshot{
		SubmitTotal:   atomic.LoadInt64(&p.submitTotal),
		SubmitSuccess: atomic.LoadInt64(&p.submitSuccess),
	}
}
