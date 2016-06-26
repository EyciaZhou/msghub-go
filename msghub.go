package main

import (
	"github.com/EyciaZhou/msghub.go/ErrorUtiles"
	"github.com/EyciaZhou/msghub.go/interface"
	_ "github.com/EyciaZhou/msghub.go/plugins/netease_news"
	_ "github.com/EyciaZhou/msghub.go/plugins/rss"
	"github.com/EyciaZhou/msghub.go/plugins/task"
	_ "github.com/EyciaZhou/msghub.go/plugins/weibo"
	"github.com/Sirupsen/logrus"
	"time"
)

func main() {
	ErrorUtiles.OUTPUT_STACK_ON_ERROR = true
	logrus.SetLevel(logrus.DebugLevel)

	Interface.Init()
	go task.TaskWatchLoop()

	for {
		time.Sleep(time.Second)
	}
}
