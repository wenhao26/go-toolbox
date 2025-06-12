package importer

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"toolbox/2025/maxmind-geoip/utils"
)

// GeoRecord 地理记录结构体
type GeoRecord struct {
	StartIP   uint32
	EndIP     uint32
	GeoID     string
	Country   string
	Region    string
	City      string
	Latitude  string
	Longitude string
}

// importGeoData 导入地理数据
func ImportGeoData(db *sql.DB, blockPath, locationPath string) error {
	// 加载位置信息
	locationMap := map[string][]string{}
	localFile, _ := os.Open(locationPath)
	localCsv := csv.NewReader(localFile)
	_, _ = localCsv.Read() // 跳过标题

	for {
		row, err := localCsv.Read()
		if err != nil {
			break
		}
		locationMap[row[0]] = []string{row[4], row[6], row[10]} // country, region, city
	}

	// 读取块数据
	blockFile, _ := os.Open(blockPath)
	blockCsv := csv.NewReader(blockFile)
	_, _ = blockCsv.Read() // 跳过标题

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO ip_geo (ip_start, ip_end, geoname_id, country_name, subdivision_name, city_name, latitude, longitude)VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	for {
		row, err := blockCsv.Read()
		if err != nil {
			break
		}

		ipRange := row[0]
		geoID := row[1]
		lat, lon := row[7], row[8]

		parts := strings.Split(ipRange, "/")
		ip := net.ParseIP(parts[0])
		mask, _ := strconv.Atoi(parts[1])
		start := utils.IpToUint32(ip)
		count := uint32(1) << (32 - mask)
		end := start + count - 1

		loc := locationMap[geoID]
		country, region, city := "", "", ""
		if len(loc) >= 3 {
			country, region, city = loc[0], loc[1], loc[2]
		}

		fmt.Printf("ip_start=%d,ip_end=%d" +
			",geoname_id=%s,country_name=%s" +
			",subdivision_name=%s,city_name=%s" +
			",latitude=%s,longitude=%s\n", start, end, geoID, country, region, city, lat, lon)

		_, err = stmt.Exec(start, end, geoID, country, region, city, lat, lon)
		if err != nil {
			log.Println("Insert error:", err)
		}
	}

	return tx.Commit()
}
