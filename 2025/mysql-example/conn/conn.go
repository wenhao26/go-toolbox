package conn

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Conn MySQL连接结构体
type Conn struct {
	DB *sql.DB
}

// NewConn 创建一个MySQL连接实例
func NewConn() (*Conn, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hinovel_test")
	if err != nil {
		return nil, err
	}

	return &Conn{DB: db}, nil
}
