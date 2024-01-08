package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/robfig/cron"

	"toolbox/MQ/rabbitmq"
)

func fileGetContents(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func jsonDecode(data []byte) ([]map[string]interface{}, error) {
	var v []map[string]interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func publish() {
	fmt.Println("执行发布")
	mq := rabbitmq.NewSimple("test.lscj")

	filename := "F:\\go-toolbox\\MQ\\amqp\\lscj.json"
	data, err := fileGetContents(filename)
	if err != nil {
		panic(fmt.Sprintf("文件读取失败：%v", err))
	}

	results, err := jsonDecode(data)
	if err != nil {
		panic(fmt.Sprintf("解析JSON失败：%v", err))
	}

	go func() {
		for _, result := range results {
			/*fmt.Println(
				result["t"],
				result["c"],
				result["zdf"],
				result["jlrl"],
				result["hsl"],
				result["qbjlr"],
				result["cddlr"],
				result["cddjlr"],
				result["ddlr"],
				result["ddjlr"],
				result["xdlr"],
				result["xdjlr"],
				result["sdlr"],
				result["sdjlr"],
			)*/

			data, _ := json.Marshal(result)
			mq.PublishSimple(string(data))
		}
	}()
}

func main() {
	c := cron.New()
	_ = c.AddFunc("*/2 * * * * *", publish)
	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
