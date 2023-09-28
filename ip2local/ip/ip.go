package ip

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TODO+ ！！！ 以下所有代码仅仅提供参考和学习 ！！！

type Serve struct {
	dataFile string
	dbFile   string
	rawDir   string
}

var (
	wg sync.WaitGroup
)

func rootPath() string {
	dir, _ := os.Getwd()
	return dir
}

func currentDate() string {
	return time.Now().Format("20060102")
}

func mkdir(pathname string) bool {
	result, err := pathIsExist(pathname)
	if err != nil {
		fmt.Printf("数据存储目录错误：%v\n", err)
		return false
	}

	if !result {
		err := os.MkdirAll(pathname, os.ModePerm)
		if err != nil {
			fmt.Printf("创建数据存储目录失败：%v\n", err)
			return false
		}
	}
	return true
}

func removeFile(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		return nil
	}

	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("删除文件失败：%v\n", err)
		return err
	}
	return nil
}

func pathIsExist(pathname string) (bool, error) {
	_, err := os.Stat(pathname)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求数据源失败：%v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("获取数据源失败：%v\n", err)
		return "", err
	}
	return string(body), nil
}

func fileGetContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

func filePutContents(filename string, data string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("写入数据源失败：%v\n", err)
		return err
	}

	n, _ := f.Seek(0, 2)
	_, err = f.WriteAt([]byte(data), n)
	defer f.Close()
	return nil
}

func appendWriteString(filename string, data string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("写入数据源失败：%v\n", err)
		return err
	}
	defer f.Close()

	write := bufio.NewWriter(f)
	_, _ = write.WriteString(data)

	// Flush将缓存的文件真正写入到文件中
	_ = write.Flush()

	return nil
}

func getFiles(pathname string, fileList []string) ([]string, error) {
	files, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("解析数据源目录失败：%v\n", err)
		return fileList, err
	}

	for _, file := range files {
		if file.IsDir() {
			_, _ = getFiles(pathname+"/"+file.Name(), fileList)
		} else {
			fileList = append(fileList, pathname+"/"+file.Name())
		}
	}
	return fileList, nil
}

func ip2long(ipAddress string) (uint, error) {
	p := net.ParseIP(ipAddress).To4()
	if p == nil {
		return 0, errors.New("invalid ipv4 format")
	}
	return uint(p[0])<<24 | uint(p[1])<<16 | uint(p[2])<<8 | uint(p[3]), nil
}

func long2ip(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}

func New() *Serve {
	return &Serve{
		dataFile: rootPath() + "/data/db/ipv4.txt",
		dbFile:   rootPath() + "/data/db/ip.db",
		//rawDir:   rootPath() + "/data/raw/" + currentDate(),
		rawDir: rootPath() + "/data/raw/latest",
	}
}

// 拉取原始数据
func (s *Serve) Pull() {
	if !mkdir(s.rawDir) {
		panic("数据存储目录未创建，无法进行数据拉取！")
	}

	urls := map[string]string{
		"delegated-apnic-latest":         "https://ftp.apnic.net/stats/apnic/delegated-apnic-latest",
		"delegated-arin-extended-latest": "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest",
		"delegated-afrinic-latest":       "https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-latest",
		"delegated-lacnic-latest":        "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-latest",
		"delegated-ripencc-latest":       "https://ftp.ripe.net/ripe/stats/delegated-ripencc-latest",
	}

	for name, url := range urls {
		wg.Add(1)
		go func(name, url string) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					return
				}
			}()

			// 获取数据源
			log.Println("执行拉取数据源：", url)
			content, err := httpGet(url)
			if err != nil {
				panic("获取数据源失败！")
			}

			// 写入数据源
			filename := s.rawDir + "/" + name + ".txt"
			_ = filePutContents(filename, content)
			log.Println("拉取完成！数据源地址：", filename)
		}(name, url)
	}
	wg.Wait()
}

// 生成ipv4数据并保存为本地文件
func (s *Serve) Create() {
	log.Println("解析最近一次拉取的数据源目录：", s.rawDir)

	var tmpArr []string
	fileList, err := getFiles(s.rawDir, tmpArr)
	if err != nil {
		panic("数据源目录不存在或者没有数据源文件，请执行[拉取]命令！")
	}

	// 删除已存在的ipv4文件
	if err := removeFile(s.dataFile); err != nil {
		panic("ipv4文件删除失败，不允许继续执行！")
	}

	for _, filename := range fileList {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					return
				}
			}()

			// 逐行读取文件，进行相关处理
			log.Println("读取数据源文件：", filename)
			f, err := os.Open(filename)
			if err != nil {
				panic("打开文件失败！")
			}
			defer f.Close()

			br := bufio.NewReader(f)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}

				line := string(a)
				line = strings.TrimSpace(line)

				if strings.Index(line, "#") == 0 || strings.Index(line, "ipv4") == -1 {
					continue
				}

				infoArr := strings.Split(line, "|")

				//var organization string
				var country string
				var typeFlag string
				var ip string
				var length string
				//var date string
				var status string

				/*if len(infoArr) >= 1 {
					organization = infoArr[0]
				}*/
				if len(infoArr) >= 2 {
					country = infoArr[1]
				}
				if len(infoArr) >= 3 {
					typeFlag = infoArr[2]
				}
				if len(infoArr) >= 4 {
					ip = infoArr[3]
				}
				if len(infoArr) >= 5 {
					length = infoArr[4]
				}
				/*if len(infoArr) >= 6 {
					date = infoArr[5]
				}*/
				if len(infoArr) >= 7 {
					status = infoArr[6]
				}

				if typeFlag != "ipv4" || (status != "assigned" && status != "allocated") {
					continue
				}

				ipStart, _ := ip2long(ip)
				ui64, _ := strconv.ParseUint(length, 10, 64)
				ipEnd := ipStart + uint(ui64) - 1
				ipEndStr, _ := long2ip(ipEnd)

				row := fmt.Sprintf("%s,%v,%v,%s,%s\n", country, ipStart, ipEnd, ip, ipEndStr)
				_ = appendWriteString(s.dataFile, row)
				log.Println("ipv4数据已写入：", row)
			}
		}(filename)
	}
	wg.Wait()
}

// 将ipv4文本文件的数据保存为SQLite3数据库文件
func (s *Serve) SaveDb() {
	db, err := sql.Open("sqlite3", s.dbFile)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(s.dataFile)
	if err != nil {
		panic("打开文件失败！")
	}
	defer f.Close()

	i := 0
	values := ""

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		line := string(a)
		line = strings.TrimSpace(line)

		infoArr := strings.Split(line, ",")
		country := infoArr[0]
		startIpNo := infoArr[1]
		endIpNo := infoArr[2]
		startIp := infoArr[3]
		endIp := infoArr[4]

		// 每1000条添加一次
		if i < 1000 {
			values += fmt.Sprintf("('%s','%s','%s','%s','%s'),", country, startIpNo, endIpNo, startIp, endIp)
		} else {
			values = strings.TrimRight(values, ",")
			insertSql := fmt.Sprintf("insert into ip_area(country,start_ip_no,end_ip_no,start_ip,end_ip) values %s", values)
			stmt, err := db.Prepare(insertSql)
			if err != nil {
				panic(err)
			}

			res, err := stmt.Exec()
			if err != nil {
				panic(err)
			}
			id, _ := res.LastInsertId()
			fmt.Println("批次插入LAST-ID：", id)

			// 重置
			i = 0
			values = ""
		}
		i++
	}
	log.Println("数据迁移完成！")
}

// 查询IP信息
func (s *Serve) Search(ip string) {
	ipNo, err := ip2long(ip)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", s.dbFile)
	if err != nil {
		panic(err)
	}

	//querySql := fmt.Sprintf("SELECT * FROM ip_area where %v>=start_ip_no and %v<=end_ip_no limit 1", ipNo, ipNo)
	//rows, err := db.Query(querySql)
	stmt, err := db.Prepare("SELECT * FROM ip_area where ?>=start_ip_no and ?<=end_ip_no limit 1")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var id int
	var country string
	var startIpNo int
	var endIpNo int
	var startIp string
	var endIp string
	_ = stmt.QueryRow(ipNo, ipNo).Scan(&id, &country, &startIpNo, &endIpNo, &startIp, &endIp)

	// TODO 返回的国家代码可以参考 README.md 文件，可以实现一些具体业务的判断等等
	fmt.Println(id, country, startIpNo, endIpNo, startIp, endIp)
}
