package main

import (
	"toolbox/beanstalk/bean"
)

func main() {
	b := bean.New()
	defer b.CloseBean()

	b.TubeStat()
}
