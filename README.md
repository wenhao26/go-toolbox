# Go Toolbox

Go 工具、第三方库使用示例与后端工程实践合集。

本仓库不是一个统一部署的业务系统，而是一组可分别阅读、运行的小项目。内容围绕 **消息队列、并发编程、任务调度、数据库访问和常用工具** 展开，部分示例使用小说阅读、内容发布、通知推送、订单导出与用户数据归档作为业务场景。

> 定位：学习参考与工程模板。包含模拟实现、实验代码和有数据修改副作用的脚本，不代表所有模块都经过生产验证。运行前请阅读对应源码与配置。

## 目录导航

- [快速开始](#快速开始)
- [重点子项目](#重点子项目)
- [按技术主题查找](#按技术主题查找)
- [业务场景示例](#业务场景示例)
- [环境与外部依赖](#环境与外部依赖)
- [运行与开发约定](#运行与开发约定)
- [安全注意事项](#安全注意事项)

## 快速开始

根模块名为 `toolbox`，当前 Go 版本声明见 [go.mod](<go.mod>)（Go 1.25.0）。个别子项目可能提出更高版本要求，请以各自说明为准。

```bash
git clone https://github.com/wenhao26/go-toolbox.git
cd go-toolbox
go version
go mod download
```

选择一个示例单独运行，不需要同时启动整个仓库。

### 示例一：生产者—消费者模型

从仓库根目录运行：

```bash
go run ./2026/producer-consumer/cmd/producer-consumer
```

该示例模拟数据生产与消费，可调整生产者数量、Worker 数量、队列容量和 QPS，并观察吞吐与退出流程。详细参数见 [使用说明](<2026/producer-consumer/README.md>)。

### 示例二：RabbitMQ 消费服务

```bash
cd rabbitmq-consumer-service
# 先修改 config.yaml，使用自己的测试 RabbitMQ 配置
go run . --config=config.yaml
```

需要可访问的 RabbitMQ 服务及相应队列。参考 [使用说明](<rabbitmq-consumer-service/README.md>)。

### 示例三：定时任务管理界面

```bash
cd 2026/gocron
# 先检查源码中的 Redis 地址和命令执行方式
go run main.go
```

启动后访问 `http://localhost:8989`。这是能够执行系统命令的管理工具，请仅用于可信的测试环境。参考 [使用说明](<2026/gocron/README.md>)。

## 重点子项目

下表是推荐阅读入口，不是全部目录的穷举清单。

| 子项目 | 内容 | 文档 / 入口 |
| --- | --- | --- |
| GoCron | Gin + Redis 定时任务管理，秒级调度、命令执行、任务与日志存储 | [文档](<2026/gocron/README.md>) · [入口](<2026/gocron/main.go>) |
| Producer-Consumer | 多生产者、有界队列、Worker Pool、吞吐统计与分阶段退出 | [文档](<2026/producer-consumer/README.md>) · [入口](<2026/producer-consumer/cmd/producer-consumer/main.go>) |
| Async Engine | 非阻塞任务投递、满队列丢弃策略、动态流量模拟与统计报表 | [文档](<2026/async-engine/README.md>) · [入口](<2026/async-engine/main.go>) |
| High-Perf Queue | 有界任务队列、对象复用与 Worker 调度示例 | [入口](<2026/high-perf-queue/cmd/main.go>) |
| Novel Creation | Redis 任务队列、处理中状态、Worker 池与任务恢复；生成逻辑为占位实现 | [文档](<2026/novel-creation/README.md>) · [入口](<2026/novel-creation/cmd/server/main.go>) |
| Omni Shutdown | 信号处理、清理钩子与优雅退出组件 | [目录](<2026/omni-shutdown/>) · [API 服务示例](<2026/omni-shutdown/examples/api-service/main.go>) |
| RabbitMQ Consumer Service | 配置加载、QoS、并发消费、手动确认、配置监听与性能分析示例 | [文档](<rabbitmq-consumer-service/README.md>) · [入口](<rabbitmq-consumer-service/main.go>) |
| Squirrel Dir Stat | 并发目录扫描，统计文件数量和总大小 | [文档](<2026/squirrel-dir-stat/README.md>) · [入口](<2026/squirrel-dir-stat/main.go>) |
| Video Parser | 基于 ffprobe 的本地 / 网络视频元信息解析工具 | [文档](<2026/videoparser/README.md>) · [入口](<2026/videoparser/cmd/main.go>) |
| MySQL Lib | 数据库访问接口、事务与预处理语句缓存示例 | [文档](<2025/go-mysql-lib/README.md>) |
| gRPC Lab | proto 定义、服务端与客户端调用实验 | [文档](<2025/grpc-lab/README.md>) |
| MaxMind GeoIP | GeoLite2 CSV 导入 SQLite，HTTP 查询 IP 地理信息 | [文档](<2025/maxmind-geoip/README.md>) |
| IP2Local | IP 数据拉取、建库、保存与查询工具 | [文档](<ip2local/README.md>) |

> 压测示例中的 QPS 参数和文档里的理论值不等于实测性能承诺；实际结果取决于硬件、任务耗时与运行配置。

## 按技术主题查找

### 消息队列与异步任务

| 目录 | 主要内容 |
| --- | --- |
| [MQ](<MQ/>) | RabbitMQ / AMQP、Kafka、NSQ 示例与封装 |
| [beanstalk](<beanstalk/>) | beanstalkd 生产、消费及队列操作 |
| [RabbitMQ 消费服务](<rabbitmq-consumer-service/>) | 并发消费服务骨架 |
| [RabbitMQ 示例](<2026/mq-example/>) | 消息发布、并发消费、手动确认与退出流程 |
| [simple-queue](<simple-queue/>) | 基于 channel 的进程内主题订阅与广播 |
| [httpmq_copy](<httpmq_copy/>) | 基于 LevelDB 的 HTTP 消息队列实验 |
| [Asynq 示例](<example/asynq/>) / [Machinery 示例](<example/machinery/>) | 异步任务框架使用 |
| [分布式任务调度](<2025/distributed-task-scheduler/>) | Master / Worker 节点实验 |

### 并发、限流与定时任务

| 目录 | 主要内容 |
| --- | --- |
| [并发练习](<2026/go-concurrency-exercises/>) | 并发聚合、超时取消、请求合并、限流与优先级调度 |
| [Ants 示例](<2025/ants-example/>) | 第三方 goroutine 池使用方式 |
| [并发搜索](<2025/concurrent-search/>) | 文件搜索、任务分发与结果汇总 |
| [threadpool](<threadpool/>) | 手写 Worker 池 |
| [timewheel](<timewheel/>) | 时间轮与延迟任务实验 |
| [文章定时发布](<2025/article-schedule/>) | Redis 有序集合实现定时发布模拟 |
| [Cron Jobs](<2026/cron-jobs/>) | 通过 HTTP API 管理定时任务 |
| [gocron-demo](<gocron-demo/>) / [cron 示例](<example/cron/>) | 调度库基础用法 |
| [任务分发](<2026/task-dist-example/>) / [高频定时器](<2026/fast-timer/>) | Worker 分发与 Ticker 实验 |

### 数据库、缓存与数据处理

| 目录 | 主要内容 |
| --- | --- |
| [mysql](<mysql/>) | GORM 使用、数据模型及业务脚本 |
| [xorm-lib](<xorm-lib/>) | XORM 示例 |
| [书籍管理示例](<2025/book-manager-example/>) | 数据库访问与分层设计 |
| [redis](<redis/>) | 字符串、Hash、HyperLogLog 等操作 |
| [mongodb](<mongodb/>) / [extend](<extend/>) | MongoDB 示例与数据库连接封装 |
| [gee-cache](<gee-cache/>) | LRU 缓存原理实现 |
| [export](<export/>) | 分批读取订单并导出 CSV |
| [read_big_csv](<read_big_csv/>) | CSV 串行 / 并发读取实验 |
| [InfluxDB 示例](<example/influxdb/>) | 时序数据库客户端使用 |

### Web、网络与工具

| 目录 | 主要内容 |
| --- | --- |
| [gin_frame](<gin_frame/>) | Gin 脚手架实验 |
| [chiapp](<chiapp/>) | chi 路由与模拟文章接口 |
| [melody-example](<melody-example/>) | WebSocket 聊天与文件监听示例 |
| [socket 示例](<example/socket/>) | TCP / UDP 通信 |
| [cmd](<cmd/>) | Cobra 命令行工具示例 |
| [conf](<conf/>) | INI 配置加载 |
| [logs_demo](<logs_demo/>) | Zap 日志与滚动文件输出 |
| [jwt](<jwt/>) / [AES](<AES/>) / [FPE 示例](<2025/fpe-example/>) | Token 与加密实验 |
| [fcm](<fcm/>) / [飞书机器人](<2025/feishu-bot/>) | 通知推送示例 |
| [download-file](<download-file/>) | HTTP 分块下载 |
| [watch-folder](<watch-folder/>) | 文件系统事件监听与文件处理 |
| [go-image-magick](<go-image-magick/>) | ImageMagick 图像处理 |
| [爬虫示例](<spider_runoob_example/>) | Colly HTML 内容提取 |
| [wire](<wire/>) | 编译期依赖注入 |
| [utils](<utils/>) / [my-lib](<my-lib/>) | ID 生成及其他辅助工具 |
| [gotest](<gotest/>) | 单元测试练习 |
| [example](<example/>) / [tryit](<tryit/>) | 其他库使用示例与语言实验 |

年度目录 [2025](<2025/>) 与 [2026](<2026/>) 中还包含其他独立实验，可按主题继续浏览。

## 业务场景示例

这些代码提供的是场景化实践，不应理解为完整的小说平台实现。

- **内容发布**：[文章定时发布](<2025/article-schedule/main.go>) 使用 Redis 存储发布时间，到期后执行模拟发布。
- **用户通知**：[FCM 任务脚本](<mysql/aikan-novel-cmd/fcm_push.go>) 演示推送任务和用户 Token 读取，当前入口并未完成实际批量发送流程。
- **订单导出**：[CSV 导出实现](<export/getdata/getdata.go>) 使用 channel 管道分批读取订单并写入文件。
- **日志处理**：[日志示例](<logs_demo/main.go>) 生成包含阅读及广告字段的模拟日志；[日志清理脚本](<mysql/aikan-novel-cmd/del_s2s_log.go>) 会执行数据库删除。
- **用户归档**：[归档迁移脚本](<2026/business-rewrite/migration_archived_user.go>) 根据注册、阅读、订单等条件决定是否归档，并涉及多表事务及物理删除。具体筛选逻辑以源码为准。
- **创作任务**：[Novel Creation](<2026/novel-creation/>) 演示异步任务消费与恢复机制，需自行接入真实生成服务。

## 环境与外部依赖

只准备所选示例需要的依赖，无需安装所有中间件。

| 依赖 | 使用场景 / 注意事项 |
| --- | --- |
| Go | 根模块声明 Go 1.25.0；部分子项目文档要求 1.25.4+ |
| Redis | 定时任务存储、延迟发布、异步任务队列及缓存示例 |
| MySQL | CRUD、订单导出、日志清理、用户归档；需要对应测试表结构 |
| RabbitMQ / Kafka / NSQ / beanstalkd | 相应消息队列示例 |
| MongoDB / InfluxDB | 相应数据库示例 |
| FFmpeg / ffprobe | 视频元信息解析；确保 ffprobe 可通过 PATH 找到 |
| ImageMagick 与 CGO 工具链 | imagick 图像处理模块；需匹配本机库与编译环境 |
| C 编译器 / CGO | 使用 go-sqlite3 的模块通常需要；不要对所有子项目统一关闭 CGO |
| protoc 与 Go 插件 | 修改 proto 后重新生成 gRPC / protobuf 代码时需要 |
| 外部 API 凭证 | FCM、飞书、Centrifugo 等；使用自己的测试凭证 |

## 运行与开发约定

### 1. 按子项目运行，而非统一启动

一些目录包含多个独立的 `package main` 文件，可能重复定义 `main` 或其他标识符。因此，不应默认整仓 `go build ./...` 或 `go test ./...` 能通过。

常见运行方式：

```bash
# 入口和配套文件构成一个包时
go run ./2026/producer-consumer/cmd/producer-consumer

# 同一目录包含多个独立示例时，只选择目标文件
go run 2025/ants-example/example0.go

# 对选定库包运行测试
go test ./jwt
```

若选定文件依赖同目录其他源文件，需要一并指定，或在没有重复入口的情况下按包运行。

### 2. 注意当前工作目录

不少示例使用相对路径读取配置、静态页面、CSV 或数据库文件。运行前确认源码中的路径；必要时进入子项目目录再执行。不要把根目录运行命令机械套用到所有模块。

### 3. 按范围构建与验证

优先对目标包运行测试、格式检查与构建。仓库中同时存在不同版本的 Redis、Cron 等依赖，修改共享依赖时需检查受影响的示例。

如果子目录实际包含独立模块定义，则以该模块为操作边界；否则使用根模块。不要为了运行单个示例随意执行 `go mod init` 或无差别升级全部依赖。

## 安全注意事项

- **不要直接连接生产环境。** 日志清理、用户迁移及导出脚本可能删除、迁移或读取敏感数据。先备份，并在隔离测试库验证筛选条件、事务和错误处理。
- **检查硬编码凭证。** 部分历史示例含 API Key、Token、连接串或服务器地址。不要复制这些值；若曾公开有效凭证，应撤销 / 轮换，并迁移至环境变量或不纳入版本控制的本地配置。
- **限制命令执行入口。** 定时任务管理模块可以执行系统命令，未经鉴权、授权和网络隔离，不应暴露到公网。
- **可靠队列不等于恰好一次处理。** Ack、重试与恢复可能造成重复执行；真实业务需补充幂等、重试上限、失败隔离及多实例协调。
- **优雅退出不等于零数据损失保证。** 有界队列可能阻塞或丢弃，内存任务无法跨进程恢复，超时退出也可能中断处理。根据业务选择持久化与退出策略。
- **评估历史依赖。** 例如旧版 JWT、FCM 与 AMQP 客户端，正式使用前应检查维护状态、兼容性与安全更新。
- **谨慎进行高负载实验。** 调低初始 QPS、Worker 数与队列容量，避免压测影响数据库、磁盘或外部服务。

---

建议阅读顺序：基础 channel / goroutine → Worker Pool 与背压 → 超时、取消和退出 → 消息确认与任务恢复 → 持久化调度与业务幂等。

