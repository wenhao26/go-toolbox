package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接池
var DB *sql.DB

// Connect 初始化数据库驱动并赋值给全局 DB
func Connect(driverName, dsn string) error {
	// 创建连接池对象
	pool, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池
	pool.SetMaxOpenConns(20)                  // 最多同时 N 个连接（并发上限）
	pool.SetMaxIdleConns(10)                  // 最多保留 N 个空闲连接常驻
	pool.SetConnMaxLifetime(30 * time.Minute) // 单个连接最长存活 N 分钟，防止长期坏连接
	pool.SetConnMaxIdleTime(5 * time.Minute)  // 空闲连接超时回收

	// 验证连接可用
	// 必须带超时 context，否则连不上会无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		err := pool.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 赋值给全局变量，之后全项目通过 db.DB 使用
	DB = pool

	return nil
}

// Close 程序退出时关闭连接池（优雅退出时调用，规范释放资源）
func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return
		}
	}
}
