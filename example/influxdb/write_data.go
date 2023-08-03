package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"gopkg.in/ini.v1"
)

var (
	influxClient influxdb2.Client
	writeAPI     api.WriteAPI
)

func init() {
	file, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}

	section := file.Section("")
	serverUrl := section.Key("server_url").String()
	token := section.Key("token").String()
	bucket := section.Key("bucket").String()
	org := section.Key("org").String()

	influxClient = influxdb2.NewClient(serverUrl, token)
	writeAPI = influxClient.WriteAPI(org, bucket)
}

func mockWrite() {
	w := mtRand(10, 100)
	h := mtRand(10, 100)
	speed := mtRand(1, 10)
	writeAPI.WriteRecord(fmt.Sprintf("robot,cmd=open w=%d,h=%d,speed=%d", w, h, speed))
	writeAPI.Flush()
}

func mtRand(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min+1) + min
}

func main() {
	defer influxClient.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()

		for {
			mockWrite()
			log.Println("running...")
			time.Sleep(time.Millisecond * 5)
			//time.Sleep(1e9)
		}
	}()
	wg.Wait()

	fmt.Println("OK")
}
