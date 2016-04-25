package main

import (
	"github.com/EyciaZhou/msghub.go/Utiles"
	"github.com/EyciaZhou/msghub.go/generant"
	_ "github.com/EyciaZhou/msghub.go/generant/netease_news"
	_ "github.com/EyciaZhou/msghub.go/generant/weibo"
	"github.com/Sirupsen/logrus"
	"time"
)

func main() {
	Utiles.OUTPUT_STACK_ON_ERROR = true
	logrus.SetLevel(logrus.DebugLevel)

	generant.Init()

	for {
		time.Sleep(time.Second)
	}
}
