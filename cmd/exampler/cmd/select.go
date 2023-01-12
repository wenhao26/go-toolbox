package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "select的使用",
	Long:  "select的使用",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("select的使用")
		doSelect()
	},
}

func init() {
	baseCmd.AddCommand(selectCmd)
}

func doSelect() {
	messageCh := make(chan string, 100)
	defer close(messageCh)

	pub(messageCh)
	sub(messageCh)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}

func pub(msgCh chan<- string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			msg := time.Now().String()
			msgCh <- msg
			fmt.Printf("[write]=%s\n", msg)
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func sub(msgCh <-chan string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			select {
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				fmt.Printf("[read]=%s\n", msg)
			}
		}
	}()
}
