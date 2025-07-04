package db_executor

import (
	"errors"
)

// - 定义常见的数据库操作错误，方便上层业务逻辑定位错误和处理

var (
	ErrNoRows               = errors.New("no rows in result set")        // 查询没有返回任何行（sql.ErrNoRows）
	ErrAffectedRowsMismatch = errors.New("affected rows mismatch")       // DML 操作影响的行数与预期不符
	ErrNoLastInsertID       = errors.New("could not get last insert ID") // 插入操作未能获取到最后插入ID
)

// IsNoRowsError 辅助函数，用于检查给定的错误是否为 ErrNoRows
func IsNoRowsError(err error) bool {
	return errors.Is(err, ErrNoRows)
}
