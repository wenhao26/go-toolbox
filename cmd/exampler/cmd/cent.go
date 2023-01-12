package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/centrifugal/gocent/v3"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var centCmd = &cobra.Command{
	Use:   "cent-publish",
	Short: "centrifugo-发布测试例子",
	Long:  "centrifugo-发布测试例子",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("centrifugo-发布")
		doCent()
	},
}

func init() {
	baseCmd.AddCommand(centCmd)
}

// 发布数据结构体
type PublishData struct {
	PubDate string `json:"pub_date"`
}

func doCent() {
	cfg, err := ini.Load("../../conf/ini/my.ini")
	if err != nil {
		panic(fmt.Sprintf("配置文件加载失败：%v", err))
	}

	sec := cfg.Section("cent")
	addr := sec.Key("addr").String()
	key := sec.Key("key").String()
	channel := sec.Key("channel").String()

	client := gocent.New(gocent.Config{
		Addr: addr,
		Key:  key,
	})

	/*pubData, _ := json.Marshal(PublishData{
		PubDate: time.Now().String(),
	})
	pubResult, err := client.Publish(context.Background(), channel, pubData)
	if err != nil {
		panic(fmt.Sprintf("调用发布时出错：%v", err))
	}
	fmt.Printf("发布到频道 %s 成功, 流位置 {offset: %d, epoch: %s}", channel, pubResult.Offset, pubResult.Epoch)*/

	// 模拟定时发布效果
	c := cron.New()
	_ = c.AddFunc("*/2 * * * * *", func() {
		testPubTimer(client, channel)
	})
	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}

}

func testPubTimer(client *gocent.Client, channel string) {
	pubData, _ := json.Marshal(PublishData{
		PubDate: time.Now().String(),
	})
	pubResult, err := client.Publish(context.Background(), channel, pubData)
	if err != nil {
		panic(fmt.Sprintf("调用发布时出错：%v", err))
	}
	fmt.Printf("发布到频道 %s 成功, 流位置 {offset: %d, epoch: %s} \n", channel, pubResult.Offset, pubResult.Epoch)
}
