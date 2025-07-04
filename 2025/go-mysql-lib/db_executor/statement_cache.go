package db_executor

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
)

// - SQL 预处理语句缓存

// StmtCache 用于缓存预处理的 SQL 语句（*sql.Stmt）
// 它是并发安全的，可以在多个协程中安全使用
type StmtCache struct {
	sync.RWMutex                      // 读写互斥锁，支持并发读写
	cache        map[string]*sql.Stmt // 存储 SQL 到预处理语句映射集
	dbClient     DBClient             // 底层数据库客户端
	logger       *log.Logger          // 日志器
}

// NewStmtCache 创建一个新的 StmtCache 实例
func NewStmtCache(dbClient DBClient, logger *log.Logger) *StmtCache {
	return &StmtCache{
		cache:    make(map[string]*sql.Stmt),
		dbClient: dbClient,
		logger:   logger,
	}
}

// GetStmt 从缓存中获取预处理语句
func (sc *StmtCache) GetStmt(ctx context.Context, query string) (*sql.Stmt, error) {
	// 尝试使用读锁获取（无阻塞，高并发）
	sc.RLock()
	stmt, ok := sc.cache[query]
	sc.RUnlock()
	if ok { // 缓存命中，直接返回
		sc.logger.Printf("[DEBUG] StmtCache: Cache hit for query '%s'", query)
		return stmt, nil
	}

	// 如果缓存中没有，则需要准备语句
	// 获取写锁，确保同一时间只有一个协程准备相同的语句，避免重复操作
	sc.Lock()
	defer sc.Unlock()

	// 双重检查锁定
	// 再次检查缓存，因为在等待写锁期间，其他协程有可能已经准备并缓存了该语句
	stmt, ok = sc.cache[query]
	if ok {
		sc.logger.Printf("[DEBUG] StmtCache: Cache hit (after double-check) for query '%s'", query)
		return stmt, nil
	}

	// 执行 PrepareContext 准备语句
	sc.logger.Printf("[DEBUG] StmtCache: Preparing new statement for query '%s'", query)
	preparedStmt, err := sc.dbClient.PrepareContext(ctx, query)
	if err != nil {
		sc.logger.Printf("[ERROR] StmtCache: Failed to prepare statement '%s': %v", query, err)
		return nil, fmt.Errorf("failed to prepare SQL statement: %w", err)
	}

	sc.cache[query] = preparedStmt
	return preparedStmt, nil
}

func (sc *StmtCache) CloseAll() {
	sc.Lock()
	defer sc.Unlock()

	sc.logger.Print("[DEBUG] StmtCache: Closing all cached statements...")
	for query, stmt := range sc.cache {
		if err := stmt.Close(); err != nil {
			sc.logger.Printf("[WARN] StmtCache: Failed to close statement '%s': %v", query, err)
		}
	}

	sc.cache = make(map[string]*sql.Stmt) // 清空缓存
	sc.logger.Print("[DEBUG] StmtCache: All cached statements closed.")
}
