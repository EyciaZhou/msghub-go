package main

import (
	_ "git.eycia.me/eycia/msghub/generant/netease_news"
	"git.eycia.me/eycia/msghub/generant"
	"time"
)

func main() {
	generant.Init()

	for {
		time.Sleep(time.Second)
	}
}
