## Geo-CSV
通过网盘分享的文件：geo-csv
链接: https://pan.baidu.com/s/1vgHRY8JfMVhnBeQ911EbUw 提取码: 5s8e


## ip_geo SQL 
CREATE TABLE IF NOT EXISTS ip_geo (
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

CREATE INDEX IF NOT EXISTS idx_ip ON ip_geo ( "ip_start" DESC, "ip_end" ASC );