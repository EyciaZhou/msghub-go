package main

import (
	"github.com/EyciaZhou/msghub.go/ErrorUtiles"
	"github.com/EyciaZhou/msghub.go/interface"
	"github.com/EyciaZhou/msghub.go/plugins"
	_ "github.com/EyciaZhou/msghub.go/plugins/netease_news"
	_ "github.com/EyciaZhou/msghub.go/plugins/rss"
	_ "github.com/EyciaZhou/msghub.go/plugins/weibo"
	"github.com/Sirupsen/logrus"
	"time"
)

func main() {
	ErrorUtiles.OUTPUT_STACK_ON_ERROR = true
	logrus.SetLevel(logrus.DebugLevel)

	Interface.Init()
	plugins.Init()

	for {
		time.Sleep(time.Second)
	}
}
