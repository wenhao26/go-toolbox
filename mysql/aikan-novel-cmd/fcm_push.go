package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NaySoftware/go-fcm"

	"toolbox/mysql/aikan-novel-cmd/global"
	"toolbox/mysql/aikan-novel-cmd/initialize"
	"toolbox/mysql/aikan-novel-cmd/storage"
)

type Token struct {
	Id       int    `json:"id"`
	UserId   int    `json:"user_id"`
	Platform int    `json:"platform"`
	Token    string `json:"token"`
}

type Task struct {
	Id        int    `json:"id"`
	BookId    int    `json:"book_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	EventType string `json:"event_type"`
	Status    int    `json:"status"`
	PushTime  int64  `json:"push_time"`
}

func (Token) TableName() string {
	return "iw_user_fcm_token"
}

func (Task) TableName() string {
	return "iw_fcm_task_push"
}

func init() {
	initialize.InitConfig()
	storage.InitDB()
}

func pullTasks() []Task {
	var task Task
	var tasks []Task

	global.DB.Debug().
		Table(task.TableName()).
		Where("status=?", 1).
		Where("is_delete=?", 0).
		Find(&tasks)
	return tasks
}

func pullTokens() []Token {
	var token Token
	var tokens []Token

	global.DB.Debug().
		Table(token.TableName()).
		Select("user_id,token").
		Order("id").
		Limit(1000).
		Find(&tokens)
	return tokens
}

func isValidTask(pushTime int64) bool {
	t := time.Now().Unix()
	return t > pushTime && t-pushTime <= 600
}

func send() {
	fcmKey := ""
	_ = ""

	data := map[string]string{
		"title": "title",
		"body":  "body",
		"msg":   "msg",
	}
	ids := []string{
		"",
	}

	c := fcm.NewFcmClient(fcmKey)
	c.NewFcmRegIdsMsg(ids, data)
	status, err := c.Send()
	if err != nil {
		panic(err)
	}

	status.PrintResults()
}

func main() {
	tasks := pullTasks()
	if len(tasks) == 0 {
		log.Println("暂无可执行的推送任务！")
		os.Exit(0)
	}

	var tokenList []string
	tokens := pullTokens()
	for _, token := range tokens {
		tokenList = append(tokenList, token.Token)
	}
	fmt.Println(tokenList)
	os.Exit(0)

	for _, task := range tasks {
		ret := isValidTask(task.PushTime)
		fmt.Println(ret)
		fmt.Println(task.Id, task.Title, task.PushTime)

	}
}

//func delLogs() {
//	var pushLog FcmPushLog
//
//	pushLogs := []FcmPushLog{}
//	db.Debug().Table(pushLog.TableName()).Select("id").Limit(delLimit).Find(&pushLogs)
//	if len(pushLogs) == 0 {
//		fmt.Println("无处理数据，退出程序")
//		os.Exit(0)
//	}
//	db.Table(pushLog.TableName()).Delete(&pushLogs)
//	log.Println("删除日志数量：", len(pushLogs))
//}
//
//func main() {
//	p, _ := ants.NewPool(10, ants.WithPreAlloc(false))
//	defer p.Release()
//
//	c := cron.New()
//	_ = c.AddFunc(fmt.Sprintf("*/%d * * * * *", cronSpec), func() {
//		for i := 0; i < runNumber; i++ {
//			_ = p.Submit(func() {
//				delLogs()
//			})
//		}
//		//for i := 0; i < runNumber; i++ {
//		//	go delLogs()
//		//}
//	})
//	c.Start()
//
//	t1 := time.NewTimer(time.Second * 10)
//	for {
//		select {
//		case <-t1.C:
//			t1.Reset(time.Second * 10)
//		}
//	}
//}
