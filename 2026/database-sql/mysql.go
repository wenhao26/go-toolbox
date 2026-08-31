package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"toolbox/2026/database-sql/db"
	"toolbox/2026/database-sql/model"
)

const (
	driverName = "mysql"
	dsn        = "root:123456@tcp(127.0.0.1:3306)/db_test?charset=utf8mb4&parseTime=true"
)

func main() {
	if err := db.Connect(driverName, dsn); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	// --------- 👇👇 使用示例 ---------
	var test model.Test

	// ✅查询单行
	err := db.DB.QueryRow("select id,title from nb_test where id = ?", 16).Scan(&test.ID, &test.Title)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(test)
	fmt.Println("---------")

	// ✅查询多行
	rows, err := db.DB.QueryContext(context.Background(), "select * from nb_test order by id limit 20")
	if err != nil {
		log.Fatal(err)
	}
	// 必须：否则连接无法归还连接池
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var tests []model.Test
	for rows.Next() {
		var t model.Test
		if err := rows.Scan(&t.ID, &t.Title, &t.CreateTime); err != nil {
			log.Fatal(err)
		}
		tests = append(tests, t)
	}

	// 循环结束后检查是否中途出错
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, t := range tests {
		fmt.Println(t)
	}
	fmt.Println("---------")

	// ✅插入、更新、删除
	res, err := db.DB.ExecContext(
		context.Background(),
		"insert into nb_test(title, create_time) values(?, ?)",
		"测试002", time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()

	fmt.Println("Insert ID: ", id)

}
