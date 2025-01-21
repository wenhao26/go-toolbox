package main

import (
	"fmt"
	"log"
	"time"

	"toolbox/2025/mysql-example/conn"
)

func main() {
	startTime := time.Now()

	// 连接MySQL
	connection, err := conn.NewConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer connection.DB.Close()

	err = connection.DB.Ping()
	if err != nil {
		panic(err)
	} else {
		log.Println("连接成功")
	}

	// 初始化游标和批次大小
	var total = 0
	var lastID = 0
	batchSize := 1000

	// 循环批次查询和处理数据
	for {
		// 查询一批数据
		rows, err := connection.DB.Query(
			`select id,sdate,book_id,author_id,read_num from hi_book_divide_data where id > ? order by id limit ?`,
			lastID,
			batchSize,
		)
		if err != nil {
			log.Fatalf("Failed to query data: %v", err)
		}

		// 记录是否还有数据
		var hasData bool = false

		// 遍历当前批次的数据
		for rows.Next() {
			hasData = true
			var bookDivideData conn.BookDivideData
			err := rows.Scan(
				&bookDivideData.ID,
				&bookDivideData.Sdate,
				&bookDivideData.BookID,
				&bookDivideData.AuthorID,
				&bookDivideData.ReadNum,
			)
			if err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}

			// TODO+
			fmt.Println(bookDivideData)
			total++

			// 更新游标位置
			lastID = bookDivideData.ID
		}

		// 关闭当前批次的查询结果
		_ = rows.Close()

		// 如果没有数据了，退出循环
		if !hasData {
			break
		}

		fmt.Println("------------------------------")
	}

	fmt.Printf("执行完毕！总行数：%d，耗时：%v\n", total, time.Since(startTime))
}
