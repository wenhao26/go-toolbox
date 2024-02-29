package main

import (
	"github.com/NaySoftware/go-fcm"
)

var (
	apiKey   = "AAAAJvC5qwU:APA91bHM1Mqqp43E_oZHEx_KwLY3F6Nsv1CqxVIw1TOsezmadsL4MsLouEp0LSRVVNaiBlcOEzGbNrPrNNWtVrhIegtFg4csmfiCLZc9oKRC1oo3lMeSR9wjPbDJEaP7w1ZZJ_IxvldB"
	senderId = "167247457029"
)

func main() {
	userToken := "fyfxdVORTfSh4ZoYfkh50r:APA91bEtVw1REVQcl8xg78aTk4i_t60FxuaNH_UvRw9UXVSG9NgX6QNoXrkBmPUx5ZWopYhb1e0M1v2u5Yng05LZxFtv3uu30OdK_mu-igzRtGejtwaSlxHGL5FzZXcF1SpySKiHFM6d"

	data := map[string]string{
		"title": "TEST-TITLE",
		"body":  "TEST-BODY",
		"msg":   "TEST_MESSAGE",
	}
	idList := []string{
		userToken,
	}

	fcmClient := fcm.NewFcmClient(apiKey)
	fcmClient.NewFcmRegIdsMsg(idList, data)
	fcmResponseStatus, err := fcmClient.Send()
	if err != nil {
		panic(err)
	}
	fcmResponseStatus.PrintResults()
}
