# Async Engine

一个用于验证高并发场景下异步任务处理可行性的实战示例。通过模拟 30000 QPS 的任务投递，验证有界队列 + Worker Pool + 背压丢弃 + 优雅退出的完整链路。

## 设计目标

- **生产与消费解耦**：投递方不受消费方处理速度限制，队列满时立即丢弃，零阻塞。
- **可控并发**：固定数量 Worker 协程消费，避免高并发下 goroutine 爆炸。
- **动态流量**：支持运行时调整目标 QPS，模拟真实业务中的流量波动。
- **优雅退出**：先停生产、再停消费，确保队列中已投递任务被处理完毕后再退出。
- **可观测**：实时看板输出 + 最终汇总报表，量化系统吞吐与背压行为。

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                        进程生命周期                          │
│                                                             │
│   [启动]              [运行]                    [退出]       │
│                                                             │
│   Producer ──Submit()──► Engine.Queue ──► Worker Pool       │
│   (100 goroutine)       (有界Channel)      (6 goroutine)    │
│        │                      │                  │          │
│        │                      │                  ▼          │
│        │                      │            模拟日志处理      │
│        │                      │            (200μs/任务)     │
│        │                      ▼                             │
│        │               队列满时丢弃                         │
│        │               (非阻塞投递)                        │
│        │                                                    │
│        ▼                                                    │
│   Limiter(令牌桶)                                           │
│   动态控制产出速率                                           │
│                                                             │
│   Reporter ◄── 每5秒采样 ───┘                               │
│   (实时看板 + 最终报表)                                      │
└─────────────────────────────────────────────────────────────┘
```

## 项目结构

```
async-engine/
├── go.mod
├── main.go              # 程序入口：生命周期编排与优雅退出
├── model/
│   └── task.go          # 任务数据模型
├── engine/
│   ├── config.go        # 引擎配置
│   ├── engine.go        # 有界队列 + 指标统计 + 优雅关闭
│   └── worker.go        # 消费协程：panic recover + 模拟处理耗时
├── producer/
│   ├── config.go        # 生产者配置与流量阶段定义
│   ├── limiter.go       # 原子令牌桶：动态速率限制
│   └── producer.go      # 高并发投递 + 阶段调度 + 背压退让
└── stats/
    └── reporter.go      # 实时看板 + 最终汇总报表
```

## 核心参数

| 组件 | 参数 | 值 | 说明 |
|------|------|-----|------|
| Engine | WorkerNum | 6 | 消费协程数 |
| Engine | QueueSize | 50000 | 队列缓冲长度 |
| Producer | WorkerNum | 100 | 投递协程数 |
| Reporter | Interval | 5s | 采样间隔 |
| 阶段1 | 5min | 10000 QPS | 稳态基线 |
| 阶段2 | 3min | 30000 QPS | 目标峰值 |
| 阶段3 | 5min | 10000 QPS | 恢复消化 |
| 阶段4 | 2min | 50000 QPS | 超限压测，触发背压 |
| 阶段5 | 15min | 30000 QPS | 持续峰值 |

单 Worker 平均处理耗时约 200 微秒，理论吞吐约 30000/s，刚好覆盖目标峰值。阶段 4 的 50000 QPS 将超过系统处理能力，触发队列背压丢弃。

## 快速开始

```bash
# 克隆项目
cd async-engine

# 编译
go build -o async-engine .

# 运行
./async-engine
```

## 运行输出

### 实时看板

```
[realtime] 10:14:32 | qps=  9998 success=     49990 drop=     0 consume=     49990 queue=     0/50000  latency=0.198ms
[realtime] 10:14:37 | qps= 30001 success=    199995 drop=     0 consume=    199995 queue=     0/50000  latency=0.201ms
[realtime] 10:14:42 | qps= 50002 success=    450010 drop= 12450 consume=    437560 queue= 50000/50000 latency=0.203ms
```

### 最终报表

```
╔════════════════════════════════════════════════════════════════════╗
║                      RUNTIME FINAL REPORT                          ║
╠════════════════════════════════════════════════════════════════════╣
║  Total Duration      : 30m0s                                       ║
║  Total Produced      : 72000000                                    ║
║  Total Success       : 64800000                                    ║
║  Total Dropped       : 7200000                                     ║
║  Total Consumed      : 64800000                                    ║
║  Max Queue Depth     : 50000                                       ║
║  Avg Consume Latency : 0.200                                       ║
║  Drop Phases         : 24                                          ║
║  Success Rate        : 90.00%                                      ║
╚════════════════════════════════════════════════════════════════════╝
```

## 关键设计决策

### 为什么用有界队列 + 非阻塞投递？

投递方调用 `Submit()` 时，如果队列已满，`select` 的 `default` 分支直接返回 `false`，不阻塞调用方。这保证了即使消费侧出现延迟或故障，生产侧也不会被拖垮。

### 为什么 Worker 需要 panic recover？

Golang 的 goroutine 没有父子关系，单个 Worker 的 panic 如果不被捕获，会导致该协程永久退出，Worker Pool 容量减少。`defer recover()` 确保单个任务的处理异常不会导致整个 Worker 死亡。

### 为什么先停 Producer 再停 Engine？

如果先 `close(engine.queue)`，而 Producer 还在投递，会触发 `panic: send on closed channel`。优雅退出的顺序必须是：**停止生产 → 关闭队列 → 等待消费完毕**。

### 令牌桶为什么用原子操作而不是 channel？

动态调整速率时，channel 方案需要重建或同步修改容量，实现复杂。原子 CAS 实现无锁、无阻塞，且支持运行时随时 `SetRate()` 调整。

## 可扩展方向

- **持久化缓冲**：队列满时不丢弃，而是写入本地 WAL 或推送到 Kafka，避免数据丢失。
- **自适应 Worker**：根据队列深度动态扩缩 Worker 数量，而非固定 6 个。
- **任务优先级**：为不同业务线分配独立队列，避免低优先级任务挤占高优先级。
- **分布式部署**：Producer 与 Engine 分离，通过 RPC 投递，支持水平扩展消费侧。

## License

MIT
