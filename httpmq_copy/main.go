package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"runtime"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const VERSION = "0.5"

var defaultMaxqueue, cacheSize, writeBuffer, cpu, keepalive *int
var ip, port, defaultAuth, dbPath *string
var db *leveldb.DB

// httpmq读取元数据api
// 从leveldb检索
// name.maxqueue-最大队列
// 名称.putpos-putpos
// 名称.getpos-getpos
func httpmqReadMetadata(name string) []string {
	maxqueue := name + ".maxqueue"
	data1, _ := db.Get([]byte(maxqueue), nil)
	if len(data1) == 0 {
		data1 = []byte(strconv.Itoa(*defaultMaxqueue))
	}

	putpos := name + ".putpos"
	data2, _ := db.Get([]byte(putpos), nil)

	getpos := name + ".getpos"
	data3, _ := db.Get([]byte(getpos), nil)

	return []string{string(data1), string(data2), string(data3)}
}

// httpmq实时getpos api
// 获取请求的httpmq当前的getpos
func httpmqNowGetpos(name string) string {
	metadata := httpmqReadMetadata(name)
	maxqueue, _ := strconv.Atoi(metadata[0])
	putpos, _ := strconv.Atoi(metadata[1])
	getpos, _ := strconv.Atoi(metadata[2])

	if getpos == 0 && putpos > 0 {
		getpos = 1 // first get operation, set getpos 1
	} else if getpos < putpos {
		getpos++ // 1nd lap, increase getpos
	} else if getpos > putpos && getpos < maxqueue {
		getpos++ // 2nd lap
	} else if getpos > putpos && getpos == maxqueue {
		getpos = 1 // 2nd first operation, set getpos 1
	} else {
		return "0" // all data in queue has been get
	}

	data := strconv.Itoa(getpos)
	_ = db.Put([]byte(name+".getpos"), []byte(data), nil)
	return data
}

// httpmq now putpos api
// get the current putpos of httpmq for request
func httpmqNowPutpos(name string) string {
	metadata := httpmqReadMetadata(name)
	maxqueue, _ := strconv.Atoi(metadata[0])
	putpos, _ := strconv.Atoi(metadata[1])
	getpos, _ := strconv.Atoi(metadata[2])

	putpos++              // increase put queue pos
	if putpos == getpos { // queue is full
		return "0" // return 0 to reject put operation
	} else if getpos <= 1 && putpos > maxqueue { // get operation less than 1
		return "0" // and queue is full, just reject it
	} else if putpos > maxqueue { //  2nd lap
		metadata[1] = "1" // reset putpos as 1 and write to leveldb
	} else { // 1nd lap, convert int to string and write to leveldb
		metadata[1] = strconv.Itoa(putpos)
	}

	db.Put([]byte(name+".putpos"), []byte(metadata[1]), nil)

	return metadata[1]
}

func init() {
	defaultMaxqueue = flag.Int("maxqueue", 1000000, "最大队列长度")
	ip = flag.String("ip", "0.0.0.0", "监听的ip地址")
	port = flag.String("port", "12138", "监听的端口")
	defaultAuth = flag.String("auth", "", "访问httpmq的auth密码")
	dbPath = flag.String("db", "level.db", "数据库路径")
	cacheSize = flag.Int("cache", 64, "缓存大小（MB）")
	writeBuffer = flag.Int("buffer", 32, "写缓冲区（MB）")
	cpu = flag.Int("cpu", runtime.NumCPU(), "httpmq的cpu数量")
	keepalive = flag.Int("k", 60, "httpmq的keepalive超时时间")
	flag.Parse()

	var err error
	db, err = leveldb.OpenFile(*dbPath, &opt.Options{BlockCacheCapacity: *cacheSize, WriteBuffer: *writeBuffer * 1024 * 1024})
	if err != nil {
		log.Fatalln("db.Get(),err:", err)
	}
}

func main() {
	runtime.GOMAXPROCS(*cpu)

	sync := &opt.WriteOptions{Sync: true}

	putnamechan := make(chan string, 100)
	putposchan := make(chan string, 100)
	getnamechan := make(chan string, 100)
	getposchan := make(chan string, 100)

	go func(chan string, chan string) {
		for {
			name := <-putnamechan
			putpos := httpmqNowPutpos(name)
			putposchan <- putpos
		}
	}(putnamechan, putposchan)

	go func(chan string, chan string) {
		for {
			name := <-getnamechan
			getpos := httpmqNowGetpos(name)
			getposchan <- getpos
		}
	}(getnamechan, getposchan)

	m := &http.ServeMux{}
	m.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		var data string
		var buf []byte
		auth := string(r.FormValue("auth"))
		name := string(r.FormValue("name"))
		opt := string(r.FormValue("opt"))
		pos := string(r.FormValue("pos"))
		num := string(r.FormValue("num"))
		charset := string(r.FormValue("charset"))

		if *defaultAuth != "" && *defaultAuth != auth {
			rw.Write([]byte("HTTPMQ_AUTH_FAILED"))
			return
		}

		method := string(r.Method)
		if method == "GET" {
			data = string(r.FormValue("data"))
		} else if method == "POST" {
			if string(r.Header.Get("Content-Type")) == "application/x-www-form-urlencoded" {
				data = string(r.FormValue("data"))
			} else {
				buf, _ = ioutil.ReadAll(r.Body)
				defer r.Body.Close()
			}
		}

		if len(name) == 0 || len(opt) == 0 {
			rw.Write([]byte("HTTPMQ_ERROR"))
			return
		}

		rw.Header().Set("Connection", "keep-alive")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Content-type", "text/plain")
		if len(charset) > 0 {
			rw.Header().Set("Content-type", "text/plain; charset="+charset)
		}

		if opt == "put" {
			if len(data) == 0 && len(buf) == 0 {
				rw.Write([]byte("HTTPMQ_PUT_ERROR"))
				return
			}

			putnamechan <- name
			putpos := <-putposchan

			if putpos != "0" {
				queueName := name + putpos
				if data != "" {
					db.Put([]byte(queueName), []byte(data), nil)
				} else if len(buf) > 0 {
					db.Put([]byte(queueName), buf, nil)
				}
				rw.Header().Set("Pos", putpos)
				rw.Write([]byte("HTTPMQ_PUT_OK"))
			} else {
				rw.Write([]byte("HTTPMQ_PUT_END"))
			}
		} else if opt == "get" {
			getnamechan <- name
			getpos := <-getposchan

			if getpos == "0" {
				rw.Write([]byte("HTTPMQ_GET_END"))
			} else {
				queueName := name + getpos
				v, err := db.Get([]byte(queueName), nil)
				if err == nil {
					rw.Header().Set("Pos", getpos)
					rw.Write(v)
				} else {
					rw.Write([]byte("HTTPMQ_GET_ERROR"))
				}
			}
		} else if opt == "status" {
			metadata := httpmqReadMetadata(name)
			maxqueue, _ := strconv.Atoi(metadata[0])
			putpos, _ := strconv.Atoi(metadata[1])
			getpos, _ := strconv.Atoi(metadata[2])

			var ungetnum float64
			var putTimes, getTimes string
			if putpos >= getpos {
				ungetnum = math.Abs(float64(putpos - getpos))
				putTimes = "1st lap"
				getTimes = "1st lap"
			} else if putpos < getpos {
				ungetnum = math.Abs(float64(maxqueue - getpos + putpos))
				putTimes = "2nd lap"
				getTimes = "1st lap"
			}

			buf := fmt.Sprintf("HTTP Simple Queue Service v%s\n", VERSION)
			buf += fmt.Sprintf("------------------------------\n")
			buf += fmt.Sprintf("Queue Name: %s\n", name)
			buf += fmt.Sprintf("Maximum number of queues: %d\n", maxqueue)
			buf += fmt.Sprintf("Put position of queue (%s): %d\n", putTimes, putpos)
			buf += fmt.Sprintf("Get position of queue (%s): %d\n", getTimes, getpos)
			buf += fmt.Sprintf("Number of unread queue: %g\n\n", ungetnum)

			rw.Write([]byte(buf))
		} else if opt == "view" {
			v, err := db.Get([]byte(name+pos), nil)
			if err == nil {
				rw.Write([]byte(v))
			} else {
				rw.Write([]byte("HTTPMQ_VIEW_ERROR"))
			}
		} else if opt == "reset" {
			maxqueue := strconv.Itoa(*defaultMaxqueue)
			db.Put([]byte(name+".maxqueue"), []byte(maxqueue), sync)
			db.Put([]byte(name+".putpos"), []byte("0"), sync)
			db.Put([]byte(name+".getpos"), []byte("0"), sync)
			rw.Write([]byte("HTTPMQ_RESET_OK"))
		} else if opt == "maxqueue" {
			maxqueue, _ := strconv.Atoi(num)
			if maxqueue > 0 && maxqueue <= 10000000 {
				db.Put([]byte(name+".maxqueue"), []byte(num), sync)
				rw.Write([]byte("HTTPMQ_MAXQUEUE_OK"))
			} else {
				rw.Write([]byte("HTTPMQ_MAXQUEUE_CANCLE"))
			}
		}
	})

	log.Fatal(http.ListenAndServe(*ip+":"+*port, m))
}
