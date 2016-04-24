package main

import (
	"github.com/EyciaZhou/msghub.go/generant"
	_ "github.com/EyciaZhou/msghub.go/generant/netease_news"
	_ "github.com/EyciaZhou/msghub.go/generant/weibo"
	"time"
	"github.com/Sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	generant.Init()

	for {
		time.Sleep(time.Second)
	}
}
