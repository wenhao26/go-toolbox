package db_executor

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// - 核心数据库操作执行器

// MySQLExecutor MySQL操作执行器
type MySQLExecutor struct {
	dbClient  DBClient
	stmtCache *StmtCache
	logger    *log.Logger
}

// NewMySQLExecutor 创建并初始化 MySQLExecutor 实例
func NewMySQLExecutor(dsn string, opts *Options) (*MySQLExecutor, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	// 创建底层数据库客户端
	dbClient, err := NewMySQLClient(dsn, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MySQL client: %w", err)
	}

	// 创建预处理语句缓存
	stmtCache := NewStmtCache(dbClient, opts.Logger)

	return &MySQLExecutor{
		dbClient:  dbClient,
		stmtCache: stmtCache,
		logger:    opts.Logger,
	}, nil
}

// Close 关闭 MySQLExecutor 所持有的数据库资源
func (e *MySQLExecutor) Close() error {
	e.logger.Print("[INFO] Closing MySQLExecutor resources...")
	e.stmtCache.CloseAll()
	return e.dbClient.Close()
}

// Exec 执行 DML（Data Manipulation Language） 语句
// 如 INSERT，UPDATE，DELETE
// 适用于不返回结果集，但会修改数据库状态的 SQL
func (e *MySQLExecutor) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// 尝试从缓存获取预处理语句
	stmt, err := e.stmtCache.GetStmt(ctx, query)
	if err != nil {
		e.logger.Printf("[ERROR] Failed to retrieve pre-processing statements: %v, Query: %s", err, query)
		return nil, fmt.Errorf("failed to retrieve pre-processing statements: %w", err)
	}

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		e.logger.Printf("[ERROR] SQL execution failed: %v, Query: %s, Args: %v", err, query, args)
		return nil, fmt.Errorf("SQL execution failed: %w", err)
	}
	return result, nil
}

// Insert 执行插入操作
func (e *MySQLExecutor) Insert(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := e.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		e.logger.Printf("[ERROR] INSERT - Failed to obtain the last inserted ID: %v, Query: %s", err, query)
		return 0, fmt.Errorf("%w: %v", ErrNoLastInsertID, err)
	}
	return lastInsertID, nil
}

// Update 执行更新操作
func (e *MySQLExecutor) Update(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := e.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		e.logger.Printf("[ERROR] UPDATE - Failed to retrieve the number of affected rows: %v, Query: %s", err, query)
		return 0, fmt.Errorf("failed to retrieve the number of affected rows: %w", err)
	}
	return rowsAffected, nil
}

// Delete 执行删除操作
func (e *MySQLExecutor) Delete(ctx context.Context, query string, args ...interface{}) (int64, error) {
	return e.Update(ctx, query, args...)
}

// Query 返回多行结果
func (e *MySQLExecutor) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	// 尝试从缓存获取预处理语句
	stmt, err := e.stmtCache.GetStmt(ctx, query)
	if err != nil {
		e.logger.Printf("[ERROR] Failed to retrieve pre-processing statements: %v, Query: %s", err, query)
		return nil, fmt.Errorf("failed to retrieve pre-processing statements: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		e.logger.Printf("[ERROR] QUERY - Query execution failed: %v, Query: %s, Args: %v", err, query, args)
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	return rows, nil
}

// QueryRow 返回单行结果
// scanFn 是一个回调函数，用于将查询结果扫描到目标结构体或变量中
// 这种设计避免了在类库内部使用反射进行通用扫描，从而保持了更好的性能
// 并将数据映射的灵活性留给调用方
func (e *MySQLExecutor) QueryRow(ctx context.Context, query string, scanFn func(row *sql.Row) error, args ...interface{}) error {
	// 尝试从缓存获取预处理语句
	stmt, err := e.stmtCache.GetStmt(ctx, query)
	if err != nil {
		e.logger.Printf("[ERROR] Failed to retrieve pre-processing statements: %v, Query: %s", err, query)
		return fmt.Errorf("failed to retrieve pre-processing statements: %w", err)
	}

	row := stmt.QueryRowContext(ctx, args...)

	// 调用 scanFn 进行结果扫描
	if err := scanFn(row); err != nil {
		if err == sql.ErrNoRows {
			e.logger.Printf("[INFO] Single line data query not found: Query: %s, Args: %v", query, args)
			return ErrNoRows // 返回自定义的 ErrNoRows
		}
		e.logger.Printf("[ERROR] QUERY_ROW -  Scan single line result failed: %v, Query: %s, Args: %v", err, query, args)
		return fmt.Errorf("scan single line result failed: %w", err)
	}
	return nil
}

// TransactionFunc 定义在事务中执行的操作函数
type TransactionFunc func(tx *sql.Tx) error

// WithTransaction 执行事务
func (e *MySQLExecutor) WithTransaction(ctx context.Context, txFn TransactionFunc, opts *sql.TxOptions) error {
	e.logger.Printf("[INFO] Transaction Start...")

	// 开启事务
	tx, err := e.dbClient.BeginTx(ctx, opts)
	if err != nil {
		e.logger.Printf("[ERROR] Failed to initiate transaction: %v", err)
		return fmt.Errorf("failed to initiate transaction: %w", err)
	}

	// 使用 defer 和 recover 来处理事务内的 panic，确保事务能够被正确回滚
	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				e.logger.Printf("[ERROR] Transaction rollback failed (Panic occurred): %v, Original Panic: %v", rollbackErr, r)
			} else {
				e.logger.Printf("[INFO] Transaction rollback successful (Panic occurred): %v", r)
			}
			panic(r) // 继续向上层传播 panic
		}
	}()

	// 执行事务函数中的定义操作
	err = txFn(tx)
	if err != nil {
		// 如果事务函数返回错误，则回滚事务
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			e.logger.Printf("[ERROR] Transaction operation failed and rollback failed: Original Error: %v, Rollback Error: %v", err, rollbackErr)
			return fmt.Errorf("transaction operation failed: %w, Rollback failed: %v", err, rollbackErr)
		}
		e.logger.Printf("[INFO] Rollback successful: %v", err)
		return fmt.Errorf("transaction operation failed: %w", err)
	}

	// 如果事务函数成功执行，则提交事务
	if commitErr := tx.Commit(); commitErr != nil {
		e.logger.Printf("[ERROR] Transaction submission failed: %v", commitErr)
		return fmt.Errorf("transaction submission failed: %w", commitErr)
	}

	e.logger.Print("[INFO] Transaction submitted successfully")
	return nil
}

// PingContext 检查数据库连接的活跃性
// 可以用户健康检查或确保数据库任然使用场景
func (e *MySQLExecutor) PingContext(ctx context.Context) error {
	e.logger.Print("[INFO] Performing database ping check")
	return e.dbClient.PingContext(ctx)
}
