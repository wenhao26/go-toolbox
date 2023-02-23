package main

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func main() {
	/*password := "1234"
	encrypt := utils.GenPassword(password)
	fmt.Println(encrypt)

	flag := utils.CheckPassword(encrypt, password)
	fmt.Println(flag)

	md5Result := utils.StrMD5(password)
	fmt.Println(md5Result)

	uuid := utils.GenUUID()
	fmt.Println(uuid)

	uuid1 := utils.GenShortUUID()
	fmt.Println(uuid1)*/

	/*id := utils.GenXID()
	fmt.Println(id)

	w := utils.NewWorker()
	orderId := w.GetID()
	fmt.Printf("EN_%d", orderId)*/

	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	/*for i := 0; i < 3; i++ {
		id := node.Generate()
		fmt.Println("id", id)
		fmt.Println("node:", id.Node(),
			"step:", id.Step(),
			"time:", id.Time(),
			"\n")
	}*/
	fmt.Println(node.Generate())

}
