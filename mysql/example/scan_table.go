package main

import (
	"fmt"
	"log"
	"time"

	"toolbox/mysql/gorm"
	"toolbox/mysql/model"
)

func main() {
	t := time.Now()

	err := gorm.Init()
	if err != nil {
		panic(fmt.Sprintf("初始化MySQL数据库失败。[err]=%s", err.Error()))
	}

	ft := model.FcmToken{}
	getDb := gorm.GDb.Table(ft.TableName())

	var count int
	page := 1
	limit := 500

	/*var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			ftList := []model.FcmToken{}

			log.Printf("正在获取第%d页数据...", page)
			offset := (page - 1) * limit
			getDb.Limit(limit).Offset(offset).Find(&ftList)
			if len(ftList) == 0 {
				log.Println("没有可加载的数据，遍历完毕！！！")
				break
			}

			for _, val := range ftList {
				fmt.Printf("platform=%d;user_id=%d;token=%s\n", val.Platform, val.UserId, val.Token)
				count++
			}

			page++
		}
		defer wg.Done()
	}()
	wg.Wait()*/

	for {
		ftList := []model.FcmToken{}

		log.Printf("正在获取第%d页数据...", page)
		offset := (page - 1) * limit
		getDb.Limit(limit).Offset(offset).Find(&ftList)
		if len(ftList) == 0 {
			log.Println("没有可加载的数据，遍历完毕！！！")
			break
		}

		for _, val := range ftList {
			fmt.Printf("platform=%d;user_id=%d;token=%s\n", val.Platform, val.UserId, val.Token)
			count++
		}

		page++
	}

	fmt.Printf("共完成 %d 条数据扫描，耗时 %s", count, time.Since(t).String())
}
