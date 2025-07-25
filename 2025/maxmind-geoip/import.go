package main

import (
	"database/sql"
	"fmt"

	"toolbox/2025/maxmind-geoip/importer"
	"toolbox/2025/maxmind-geoip/utils"
)

func main() {
	db, err := sql.Open("sqlite3", "./geoip.db")
	if err != nil {
		panic(err)
	}

	// 创建存储表（不存在则创建）
	err = utils.CreateTable(db)
	if err != nil {
		panic(err)
	}

	err = importer.ImportGeoData(db, "./GeoLite2-City-Blocks-IPv4.csv", "./GeoLite2-City-Locations-zh-CN.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println("导入完成")
}
