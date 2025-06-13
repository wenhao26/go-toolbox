package utils

import (
	"database/sql"
	"encoding/binary"
	"net"
)

// IpToUint32 IP转Uint32
func IpToUint32(ip net.IP) uint32 {
	ipv4 := ip.To4()
	return binary.BigEndian.Uint32(ipv4)
}

// CreateTable 创建存储表
func CreateTable(db *sql.DB) error {
	// 如果`ip_geo`表不存在，则创建
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ip_geo (
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
	return err
}
