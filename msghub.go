package main

import (
	"github.com/EyciaZhou/msghub.go/generant"
	_ "github.com/EyciaZhou/msghub.go/generant/netease_news"
	_ "github.com/EyciaZhou/msghub.go/generant/weibo"
	"time"
)

func main() {
	generant.Init()

	for {
		time.Sleep(time.Second)
	}
}
