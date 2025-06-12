package main

import (
	"database/sql"
	"fmt"

	"toolbox/2025/maxmind-geoip/importer"
)

func main() {
	db, err := sql.Open("sqlite3", "./geoip.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ip_geo (
				  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				  "ip_start" INTEGER NOT NULL DEFAULT 0,
				  "ip_end" INTEGER NOT NULL DEFAULT 0,
				  "geoname_id" TEXT NOT NULL DEFAULT '',
				  "country_name" TEXT NOT NULL DEFAULT '',
				  "subdivision_name" TEXT NOT NULL DEFAULT '',
				  "city_name" TEXT NOT NULL DEFAULT '',
				  "latitude" TEXT NOT NULL DEFAULT '',
				  "longitude" TEXT NOT NULL DEFAULT ''
				);
				CREATE INDEX IF NOT EXISTS idx_ip ON ip_geo ("ip_start" DESC, "ip_end" ASC);
			`)
	if err != nil {
		panic(err)
	}

	err = importer.ImportGeoData(db, "./GeoLite2-City-Blocks-IPv4.csv", "./GeoLite2-City-Locations-zh-CN.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println("导入完成")
}
