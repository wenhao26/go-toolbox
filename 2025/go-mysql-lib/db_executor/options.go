package db_executor

import (
	"log"
	"time"
)

// - 定义初始化数据库连接和类库所需的配置参数

// Options 配置 MySQLExecutor 行为和底层数据库连接池
type (
	Options struct {
		// 最大打开的数据库连接数，控制并发量
		// 这里直接控制同时允许多少个并发请求和数据库交互
		// 设置过高可能耗尽数据库资源，过低可能导致请求排队，需要根据实际负载和数据库性能调整
		MaxOpenConns int

		// 连接池中最大空闲连接数
		// 空闲连接保持连接活跃，当有新的请求时可以直接服用，减少新建连接开销
		// 对于高并发短连接场景尤其重要
		MaxIdleConns int

		// 连接可被复用的最长时间，防止连接过期或失效
		// 超过这个时间的连接会被关闭并重新建立，有助于防止连接长时间不刷新导致的问题
		// 例如数据库连接被断开，或负载均衡将请求路由到已关闭的连接
		ConnMaxLifetime time.Duration

		// 连接在被关闭前可保持空闲的最长时间
		// 如果一个连接在这个时间内没有被使用，它将被从连接池中移除并关闭
		// 避免占用不必要的数据库资源
		ConnMaxIdleTime time.Duration

		// 数据库连接初始化的 ping 超时时间
		// 在首次建立连接或者检查连接活跃性时，会发送一个 Ping 操作
		// 如果在这个时间内没有响应，则认为连接失败
		PingTimeout time.Duration

		// 可选的日志器，用于记录 SQL 的操作、性能统计和错误信息
		Logger *log.Logger
	}
)

// DefaultOptions 默认配置
func DefaultOptions() *Options {
	return &Options{
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
		PingTimeout:     5 * time.Second,
		Logger:          log.Default(),
	}
}
