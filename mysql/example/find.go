package main

import (
	"fmt"

	"toolbox/mysql/gorm"
	"toolbox/mysql/model"
)

func main() {
	err := gorm.Init()
	if err != nil {
		panic(fmt.Sprintf("初始化MySQL数据库失败。[err]=%s", err.Error()))
	}

	article := model.Article{}
	getDb := gorm.GDb.Table(article.TableName())

	// 获取第一条记录（主键升序）
	//getDb.Find(&article)

	// 获取最后一条记录（主键降序）
	//getDb.Last(&article)

	// 获取一条记录，没有指定排序字段
	//getDb.Take(&article)

	//result := getDb.First(&article)
	//fmt.Println(result.RowsAffected, result.Error)

	// 检查 ErrRecordNotFound 错误
	//errors.Is(result.Error, gorm2.ErrRecordNotFound)

	// 获取多条记录
	articles := []model.Article{}
	getDb.Limit(2).Find(&articles)

	fmt.Println(articles)

}
