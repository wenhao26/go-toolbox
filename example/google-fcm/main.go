package main

import (
	"fmt"

	"github.com/NaySoftware/go-fcm"
)

const (
	// 爱看FCM-KEY
	key = "AAAAJvC5qwU:APA91bHM1Mqqp43E_oZHEx_KwLY3F6Nsv1CqxVIw1TOsezmadsL4MsLouEp0LSRVVNaiBlcOEzGbNrPrNNWtVrhIegtFg4csmfiCLZc9oKRC1oo3lMeSR9wjPbDJEaP7w1ZZJ_IxvldB"
)

func main() {
	data := map[string]string{
		"title": "__TITLE__",
		"body":  "_BODY__",
		"msg":   "__MSG__",
	}
	tokens := []string{
		"fyfxdVORTfSh4ZoYfkh50r:APA91bEtVw1REVQcl8xg78aTk4i_t60FxuaNH_UvRw9UXVSG9NgX6QNoXrkBmPUx5ZWopYhb1e0M1v2u5Yng05LZxFtv3uu30OdK_mu-igzRtGejtwaSlxHGL5FzZXcF1SpySKiHFM6d",
	}

	c := fcm.NewFcmClient(key)
	c.NewFcmRegIdsMsg(tokens, data)

	c.AppendDevices(tokens)

	status, err := c.Send()
	if err != nil {
		panic(err)

	}
	status.PrintResults()
	fmt.Println(status)
}
