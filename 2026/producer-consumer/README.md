# Producer-Consumer 高并发处理示例（Go 1.25.4 / Windows）

一个可直接运行、生产级代码风格的**生产者-消费者**并发模型示例：多个生产者模拟实时数据推送，可配置数量的
worker 并发消费并模拟不同耗时的业务处理，支持优雅退出、健壮性兜底与实时吞吐监控。

---

## 一、目录结构

```
producer-consumer/
├── go.mod                                 # module 定义，go 1.25.4
├── README.md                              # 本文档
├── cmd/
│   └── producer-consumer/
│       └── main.go                        # 程序入口：初始化日志、监听退出信号、启动 App
└── internal/
    ├── app/
    │   └── app.go                         # 应用生命周期编排（启动/运行/优雅停机/强制兜底）
    ├── config/
    │   └── config.go                      # 命令行参数解析，所有可调优参数集中管理
    ├── model/
    │   └── message.go                     # 消息数据结构与消息类型定义
    ├── producer/
    │   └── producer.go                    # 生产者：按 QPS 持续生成消息，模拟实时数据推送
    ├── consumer/
    │   └── consumer.go                    # 消费者：固定大小 worker pool，模拟不同耗时处理
    └── stats/
        └── stats.go                       # 并发安全的原子计数器，实时统计吞吐/成功/失败
```

采用 Go 官方推荐的 `cmd/` + `internal/` 分层：`cmd` 只负责组装（wiring），
业务逻辑全部下沉到 `internal` 的各个职责单一的包中，包与包之间通过接口/channel 解耦，
`internal` 保证这些包不会被模块外的代码误引用。

---

## 二、快速开始（Windows）

```powershell
# 1. 进入项目目录
cd producer-consumer

# 2. 直接运行（默认 3 个生产者、8 个 worker、QPS=200/生产者，一直运行直到 Ctrl+C）
go run ./cmd/producer-consumer

# 3. 自定义参数运行
go run ./cmd/producer-consumer `
  -producer=4 `
  -worker=32 `
  -queue=2000 `
  -qps=500 `
  -stats-interval=1s `
  -shutdown-timeout=5s `
  -log-level=info

# 4. 编译为可执行文件
go build -o producer-consumer.exe ./cmd/producer-consumer
.\producer-consumer.exe -worker=16

# 5. 运行 60 秒后自动退出（用于压测/CI）
go run ./cmd/producer-consumer -duration=60s
```

按 `Ctrl+C` 即可触发优雅退出流程，日志会依次打印停机的三个阶段。

### 命令行参数一览

| 参数                | 默认值 | 说明                                             |
|---------------------|--------|--------------------------------------------------|
| `-producer`          | 3      | 并发生产者数量                                    |
| `-worker`            | 8      | 消费者 worker 数量                                |
| `-queue`             | 1000   | 生产者-消费者之间的缓冲队列容量（channel buffer） |
| `-qps`               | 200    | 单个生产者每秒生产的消息数                        |
| `-duration`          | 0      | 运行时长，0 表示一直运行直到 Ctrl+C               |
| `-shutdown-timeout`  | 10s    | 优雅退出时等待 worker 消费完剩余消息的最长时间    |
| `-stats-interval`    | 2s     | 统计信息打印周期                                  |
| `-log-level`         | info   | 日志级别：debug/info/warn/error                   |

---

## 三、架构设计思路

```
┌────────────┐     ┌────────────┐     ┌────────────┐
│ Producer 0 │──┐  │            │  ┌─▶│  Worker 0  │
├────────────┤  │  │            │  │  ├────────────┤
│ Producer 1 │──┼─▶│  channel   │──┼─▶│  Worker 1  │
├────────────┤  │  │ (有界队列) │  │  ├────────────┤
│ Producer N │──┘  │            │  └─▶│  Worker M  │
└────────────┘     └────────────┘     └────────────┘
     写入(背压)                          竞争消费(负载均衡)
```

### 1. 核心模型：channel 作为有界队列 + Fan-in / Fan-out

- **Fan-in**：多个 Producer goroutine 并发向同一个 `chan *model.Message` 写入，
  Go 的 channel 本身就是多写多读并发安全的，不需要额外加锁。
- **Fan-out**：多个 Worker goroutine 并发从同一个 channel 读取，天然实现
  "谁先处理完谁再取下一条" 的负载均衡，不需要额外的调度器/分发器。
- **有界队列即背压（Backpressure）**：`make(chan *model.Message, QueueSize)` 使用带缓冲
  channel。当消费能力跟不上生产速度、队列写满后，`Producer.Run` 里的
  `p.out <- msg` 会自然阻塞，从而把压力反向传导给生产端，防止无限制的内存增长——
  这是比"无界队列 + 内存爆炸"更健壮的工程实践。

### 2. 优雅退出：四阶段生命周期编排（`internal/app/app.go`）

这是本项目设计的核心，也是大厂线上服务停机最常用的范式：

1. **停止生产**：先 `cancel` 一个只作用于 Producer 的 `producerCtx`，
   所有生产者停止生成新消息并退出，用 `sync.WaitGroup` 等待其全部退出。
2. **关闭队列**：`close(channel)`，这是 Go 中"通知下游没有更多数据了"的标准手法。
3. **排空消费**：Worker 的主循环使用 `for msg := range channel`，
   channel 关闭后会继续把已经在队列里的消息处理完，再自然退出 `range` 循环——
   保证了**停机不丢数据**。
4. **超时强制兜底**：如果排空过程超过 `-shutdown-timeout`（例如某条消息处理逻辑卡死），
   会 `cancel` 另一个只作用于 Worker 内部处理的 `workerCtx`，
   正在 `select` 等待处理结果的 worker 会被立刻打断退出，避免进程无法关闭——
   保证了**停机不会无限等待**。

这种"两个独立 context 分别控制生产端和消费端的退出节奏"是本设计与许多简化版本
（只用一个 context 广播取消）的关键区别：广播式取消会导致队列中尚未处理的消息被直接丢弃，
而这里的分阶段设计能做到尽量不丢数据，同时仍然有兜底超时保证进程可以退出。

### 3. 健壮性设计

- **panic 双重兜底**：`worker()` 和 `process()` 分别有一层 `recover()`。
  单条消息处理 panic 只会被 `process()` 里的 recover 捕获，不影响该 worker
  继续处理下一条消息；即便更严重的 panic 逃逸到 `worker()` 层，也不会导致整个进程崩溃。
- **channel 方向类型约束**：`Producer.out` 声明为 `chan<-`（只写），
  `Pool.in` 声明为 `<-chan`（只读），编译期即可防止误用（例如消费者误把消息写回队列）。
- **不使用全局变量/单例**：所有状态（`Stats`、`channel`、`seq`）都通过依赖注入的方式
  显式传递，避免隐式全局状态带来的测试困难和并发陷阱。

### 4. 高并发性能设计

- **无锁计数器**：`internal/stats` 使用 `atomic.Uint64` 而非 `mutex + int64`，
  在高频写入（每条消息都要 `+1`）场景下避免锁竞争开销。
- **`math/rand/v2` 包级并发安全随机数**：Go 1.22 起 `math/rand/v2` 的包级函数
  默认使用并发安全的 ChaCha8 算法，避免了旧版 `math/rand` 需要每个 goroutine
  单独持有 `*rand.Rand` 或者用 `sync.Mutex` 包裹全局 `Source` 的问题。
- **`time.Ticker` 限速**：生产者用 `Ticker` 而非 `time.Sleep` 循环来控制 QPS，
  精度更高、runtime 调度开销更低，也更符合"实时推送"的语义。
- **`time.Timer` 模拟处理耗时**：消费者用 `NewTimer` + `select` 而不是直接
  `time.Sleep`，是为了能同时监听 `ctx.Done()` 实现处理过程中的可中断性（用于强制兜底）。

### 5. 可观测性

`internal/app.reportStats` 周期性输出：
`produced`（累计生产）/ `consumed`（累计消费尝试）/ `succeeded`（成功）/ `failed`（失败或被中断）/
`queue_len`、`queue_cap`（队列积压情况）/ `goroutines`（当前协程数，辅助判断是否存在泄漏）。
这些指标是压测和线上排障时最关键的第一手信息。

---

## 四、涉及的知识点清单

| 分类           | 知识点                                                                 |
|----------------|--------------------------------------------------------------------------|
| 并发原语       | goroutine、有缓冲/无缓冲 channel、`select`、channel 的关闭与 `range` 语义 |
| 并发原语       | `sync.WaitGroup`（等待一组协程结束）、单向 channel 类型（`chan<-` / `<-chan`）|
| 并发原语       | `sync/atomic`（`atomic.Uint64`）无锁计数器                                |
| 上下文控制     | `context.WithCancel`、多个独立 context 分别控制不同子系统的生命周期      |
| 优雅退出       | `signal.NotifyContext` 监听 `os.Interrupt` / `syscall.SIGTERM`           |
| 优雅退出       | "先关水龙头、再排空队列、最后超时强制退出" 的三段式停机模型               |
| 健壮性         | `defer` + `recover()` 的多层 panic 兜底                                  |
| 限流/调度      | `time.Ticker` 做生产限速、`time.Timer` 模拟可中断的耗时处理               |
| 随机数         | `math/rand/v2` 包级并发安全随机源（Go 1.22+）                            |
| 日志           | `log/slog` 结构化日志（Go 1.21+ 标准库）                                 |
| 工程结构       | `cmd/` + `internal/` 分层、依赖注入、单一职责包划分                       |
| 背压           | 有界 channel 天然实现生产端限速（backpressure）                          |
| 可观测性       | 原子计数器 + 周期性指标上报（吞吐、队列深度、goroutine 数）              |

---

## 五、扩展方向（可选）

- **消息优先级**：用多个 channel（高/中/低优先级）+ 带权重的 `select` 实现优先级队列。
- **动态扩缩容**：根据 `queue_len` 动态增加/减少 worker 数量（需要引入可取消的 worker 生命周期管理）。
- **重试与死信队列**：`process()` 失败后不直接丢弃，而是重新入队或写入死信 channel/文件。
- **接入真实中间件**：将 channel 替换为 Kafka / RabbitMQ / NATS 等消息队列 client，
  当前的 Producer/Consumer 接口设计（面向 channel）可以平滑迁移。
- **Prometheus 指标**：把 `internal/stats` 的计数器改为 `prometheus.Counter`，暴露 `/metrics`。

---

## 六、验证方式

由于当前开发环境无 Go 编译器/无网络，代码已经过人工逐行 review 确保语法与类型正确，
请在具备 Go 1.25.4（或兼容版本，如 1.22+，因为使用了 `math/rand/v2` 与内建 `slog`）
的 Windows 环境中执行：

```powershell
go vet ./...
go build ./...
go run ./cmd/producer-consumer -duration=10s -worker=16 -qps=300
```

预期现象：日志周期性打印 `stats snapshot`，`produced` 与 `consumed` 持续增长；
10 秒后自动进入优雅停机流程，依次打印 `step 1/3` ~ `step 3/3`，最终打印
`final stats` 与 `application exited cleanly, bye`。