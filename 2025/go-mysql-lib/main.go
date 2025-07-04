package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"toolbox/2025/go-mysql-lib/db_executor"
)

// User 用户数据
type User struct {
	ID         int64
	Name       string
	Email      string
	Gender     int8
	CreateTime int
}

func main() {
	// 初始化配置选项
	opts := db_executor.DefaultOptions()

	// 创建数据库执行器
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?parseTime=true&charset=utf8"
	executor, err := db_executor.NewMySQLExecutor(dsn, opts)
	if err != nil {
		log.Fatalf("初始化数据库执行器失败: %v", err)
	}
	defer executor.Close()

	// 创建一个带有超时的上下文，用于控制所有数据库操作的执行时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- 测试 ---

	// 查询单条数据
	var user User
	queryOneSql := "select id,name,email,gender,create_time from user where id = ?"
	err = executor.QueryRow(ctx, queryOneSql, func(row *sql.Row) error {
		return row.Scan(&user.ID, &user.Name, &user.Email, &user.Gender, &user.CreateTime)
	}, 46363)

	if err != nil {
		if db_executor.IsNoRowsError(err) { // 使用自定义的错误判断函数
			fmt.Println("未找到 ID 为 `46363` 的用户")
		} else {
			fmt.Printf("查询用户 `46363` 失败: %v\n", err)
		}
	} else {
		fmt.Printf("查询到用户 `46363` : %+v\n", user)
	}

}
