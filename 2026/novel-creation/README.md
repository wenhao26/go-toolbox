# Novel Creation AI Service

## 运行要求
- Go 1.25.4+
- Redis 6.0+

## 配置
复制 `config.yaml` 根据实际环境修改 Redis 地址、并发数等参数。

## 编译运行
```bash
go mod tidy
go run cmd/server/main.go

## 项目目录结构
```text
novel-creation/
├── cmd/
│   └── server/
│       └── main.go                # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go              # 配置加载
│   ├── queue/
│   │   └── redis_queue.go         # Redis队列操作封装
│   ├── worker/
│   │   └── worker.go              # Worker工作逻辑
│   └── handler/
│       └── task_handler.go        # 具体任务处理器（可扩展）
├── pkg/
│   └── logger/
│       └── logger.go              # 日志封装
├── config.yaml                    # 配置文件示例
├── go.mod
├── go.sum
└── README.md
```

