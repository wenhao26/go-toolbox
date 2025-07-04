package db_executor

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// - 底层数据库客户端接口与实现

// DBClient 定义底层数据库操作接口
type DBClient interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	PingContext(ctx context.Context) error
	Close() error
}

// mySQLClient 实现 DBClient 接口，封装 *sql.DB
type mySQLClient struct {
	db     *sql.DB     // 标准库连接池对象
	logger *log.Logger // 日志器
}

// NewMySQLClient 创建 MySQLClient 实例并初始化数据库连接池
func NewMySQLClient(dsn string, opts *Options) (DBClient, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}

	// 应用配置选项到数据库连接池
	// 这些配置对性能和高并发至关重要
	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetConnMaxLifetime(opts.ConnMaxIdleTime)
	db.SetConnMaxIdleTime(opts.ConnMaxIdleTime)

	// 使用 Context 进行 Ping 操作，并设置超时，确保数据库连接是否可用
	ctx, cancel := context.WithTimeout(context.Background(), opts.PingTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connection to database failed (Ping check failed): %w", err)
	}

	return &mySQLClient{
		db:     db,
		logger: opts.Logger,
	}, nil
}

// ExecContext 执行 DML（Data Manipulation Language） 语句
// 如 INSERT，UPDATE，DELETE
// 适用于不返回结果集，但会修改数据库状态的 SQL
func (c *mySQLClient) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	c.logger.Printf("[DEBUG] Exec: Query='%s', Args='%v'", query, args) // 记录执行的 SQL 和参数
	return c.db.ExecContext(ctx, query, args...)
}

// QueryContext 执行 SELECT 语句，返回多行结果集
func (c *mySQLClient) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	c.logger.Printf("[DEBUG] Query: Query='%s', Args='%v'", query, args) // 记录查询的 SQL 和参数
	return c.db.QueryContext(ctx, query, args...)
}

// QueryRowContext 执行 SELECT 语句，期望返回单行结果
// 如果查询返回多行，它只会读取第一行
// 如果没有行返回，Scan 方法会返回 sql.ErrNoRows
func (c *mySQLClient) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	c.logger.Printf("[DEBUG] QueryRow: Query='%s', Args='%v'", query, args) // 记录单行查询的 SQL 和参数
	return c.db.QueryRowContext(ctx, query, args...)
}

// PrepareContext 创建 SQL 预处理语句
// 预处理语句在性能上非常重要，因为它允许数据库预先解析和优化 SQL
// 后续多次执行时只需发送参数，减少了重复解析的开销
func (c *mySQLClient) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	c.logger.Printf("[DEBUG] Prepare: Query='%s'", query) // 记录预处理的 SQL
	return c.db.PrepareContext(ctx, query)
}

// BeginTx 开启一个数据库事务
// 事务用于确保一系列数据库操作的原子性（要么全部成功，要么全部失败）
func (c *mySQLClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	c.logger.Printf("[DEBUG] BeginTx: Options='%+v'", opts) // 记录事务开启
	return c.db.BeginTx(ctx, opts)
}

// PingContext 检查数据库连接的活跃性
func (c *mySQLClient) PingContext(ctx context.Context) error {
	c.logger.Print("[DEBUG] Ping database") // 记录 Ping 操作
	return c.db.PingContext(ctx)
}

// Close 关闭底层数据库连接池
// 在应用程序退出时调用，释放数据库资源
func (c *mySQLClient) Close() error {
	c.logger.Print("[DEBUG] Closing database connection") // 记录关闭操作
	return c.db.Close()
}
