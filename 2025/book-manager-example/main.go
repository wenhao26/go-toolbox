package main

import (
	"fmt"
	"log"

	"toolbox/2025/book-manager-example/config"
	"toolbox/2025/book-manager-example/model"
	"toolbox/2025/book-manager-example/repo"
	"toolbox/2025/book-manager-example/service"
)

func main() {
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建数据访问层和业务逻辑层
	bookRepo := repo.NewBookRepo()
	bookSrv := service.NewBookSrv(bookRepo)

	// 示例：添加书籍
	newBook := &model.Book{
		Title:  "1111111",
		Author: "22",
		Price:  8.88,
	}
	err := bookSrv.CreateBook(newBook)
	if err != nil {
		log.Printf("Error creating book: %v", err)
	} else {
		fmt.Println("New book created successfully")
	}

	// 示例：查询书籍
	book, err := bookSrv.GetBook(1)
	if err != nil {
		log.Printf("Error getting book: %v", err)
	} else {
		fmt.Printf("Book found: %+v\n", book)
	}
}
