package main

import (
	"encoding/json"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"fmt"
	"sync"
	"time"
)

var gDB *gorm.DB

func init() {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, "root", "root", "127.0.0.1", "3306", "istory_db", "utf8mb4")
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tp_",
			SingularTable: true,
		},
		Logger:               logger.Default.LogMode(logger.Info),
		DisableAutomaticPing: false,
	})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}

	fmt.Println("连接成功")
	fmt.Println(Db)
	gDB = Db
}

type Expend struct {
	ExpendId  int
	UserId    int
	BookId    int
	SectionId int
}

func main() {
	start := time.Now()

	// 处理行数
	count := 10215696
	// 每页处理条数
	pageSize := 500
	// 总页数
	page := (count + pageSize - 1) / pageSize

	wg := sync.WaitGroup{}

	for i := 0; i < page; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				if err := recover(); err != nil {
					return
				}
			}()
			defer wg.Done()
			// 计算当前页的偏移量
			offset := i * pageSize
			expendList := []Expend{}
			gDB.Limit(pageSize).Offset(offset).Find(&expendList)

			data, _ := json.Marshal(expendList)
			fmt.Println(string(data) + "\n")
		}(i)
	}
	wg.Wait()

	fmt.Println("--花费的时间：", time.Since(start).String())
}
