package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var chCmd = &cobra.Command{
	Use:   "channel",
	Short: "channel的使用",
	Long:  "channel的使用",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("channel的使用")
		doChannel()
	},
}

func init() {
	baseCmd.AddCommand(chCmd)
}

func doChannel() {
	// 无缓冲：模拟外卖场景，双方阻塞
	ch := make(chan string)
	go func() {
		who := "美团外卖-"
		food := "肯德基套餐"
		fmt.Println(who, "送餐中...2s")
		time.Sleep(2e9)
		fmt.Println(who, "外卖小哥到门口等候，等到客人开门...")
		ch <- food
		fmt.Println(who, "送餐完成！")
	}()
	go func() {
		who := "顾客-"
		fmt.Println(who, "等待外卖中...")
		time.Sleep(3e9)
		fmt.Println(who, "撒泡尿再去开门中...3s")
		food := <-ch
		fmt.Println(who, "拿到了", food, "外卖！")
	}()

	// 有缓冲
	/*ch := make(chan string, 1)
	go func() {
		who := "美团外卖-"
		food := "肯德基套餐"
		fmt.Println(who, "送餐中...2s")
		time.Sleep(2e9)
		fmt.Println(who, "外卖小哥已到达，把外卖放在门口了...")
		ch <- food
		fmt.Println(who, "送餐完成！继续跑单...")
	}()
	go func() {
		who := "顾客-"
		fmt.Println(who, "等待外卖中...")
		time.Sleep(3e9)
		fmt.Println(who, "让外卖小哥把外卖放在门口，撒泡尿再去取餐...3s")
		food := <-ch
		fmt.Println(who, "拿到了", food, "外卖！")
	}()*/

	time.Sleep(5e9)
}
