package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"toolbox/2025/maxmind-geoip/utils"
)

// Result 查询结果结构体
type Result struct {
	IP        string `json:"ip"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

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

	// 启动 HTTP 服务
	http.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.URL.Query().Get("address")
		ip := net.ParseIP(ipStr)
		if ip == nil {
			http.Error(w, "无效的IP", 400)
			return
		}

		ipInt := utils.IpToUint32(ip)
		row := db.QueryRow(
			`SELECT country_name, subdivision_name, city_name, latitude, longitude FROM ip_geo WHERE ip_start <= ? AND ip_end >= ? LIMIT 1`,
			ipInt,
			ipInt,
		)

		var res Result
		res.IP = ipStr
		err := row.Scan(&res.Country, &res.Region, &res.City, &res.Latitude, &res.Longitude)
		if err != nil {
			http.Error(w, "未找到IP", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	fmt.Println("Server started at :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))

	// 访问示例：http://localhost:8085/ip?address={IP地址}
}
