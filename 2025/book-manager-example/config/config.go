package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB 是数据库的全局连接池
var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	dsn := "root:password@tcp(localhost:3306)/book_db?charset=utf8&parseTime=True&loc=Local"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		return err
	}

	fmt.Println("Connected to the database successfully!")
	return nil
}
