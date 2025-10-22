package main

import (
	"fmt"
	"time"

	"toolbox/2025/feishu-bot/notifier"
)

const (
	webHook   = "xxxxx"
	secretKey = "xxxxx"
)

func main() {
	timestamp := time.Now().Unix()
	bot := notifier.NewNotifier(webHook, secretKey)

	message := notifier.Text("The message comes from golang")
	message.Payload["timestamp"] = timestamp
	message.Payload["sign"] = bot.Sign(timestamp)

	response := bot.Send(*message)
	fmt.Println(response)
}
