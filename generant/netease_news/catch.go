package netease_news

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"errors"
)

type NeteaseNewsCatchConfigure struct {
	NeteaseNewsChannel
	PagesOneTime		int
	DelayBetweenPage	time.Duration
	DelayBetweenEachCatchRound time.Duration

	quit chan int
}

func NewDefaultNeteaseNewsCatchConfigure(channelName string) (*NeteaseNewsCatchConfigure, error) {
	if hv, _ := channelsDefault[channelName]; !hv {
		return nil, errors.New("no such channel")
	}
	channelInfo := channelsDefault[channelName]
	return &NeteaseNewsCatchConfigure{
		*channelInfo,
		2,
		time.Second * 10,
		time.Minute * 10,

		make(chan int),
	}
}

func NewNeteaseNewsCatchConfigure(chann NeteaseNewsCatchConfigure, pagesOneTime int,
					delayBetweenPage, delayBetweenEachCatchRound time.Duration) {
	return &NeteaseNewsCatchConfigure{
		chann,
		pagesOneTime,
		delayBetweenPage,
		delayBetweenEachCatchRound,

		make(chan int),
	}
}

func (conf *NeteaseNewsCatchConfigure)catchOneTime() {
	log.WithFields(log.Fields{
		"channel" : conf.Name,
		"page num" : conf.PagesOneTime,
	}).Info("Start Catch")

	cnt := 0

	for i := 0; i < conf.PagesOneTime; i++ {
		if i != 0 {
			select {
			case <-time.After(conf.DelayBetweenPage) :

			case <-conf.quit :
				return
			}
		}

		log.WithFields(log.Fields{
			"channel" : conf.Name,
			"page no." : i,
		}).Info("Start Catch Page")

		news, err := getNewsList(conf.ID, i)
		if err != nil {
			log.WithFields(log.Fields{
				"channel" : conf.Name,
				"page no." : i,
				"error" : err.Error(),
			}).Error("Errors when catch")
			continue
		}
		log.WithFields(log.Fields{
			"channel" : conf.Name,
			"page no." : i,
		}).Info("Fetch Page ends")

		for _, item := range news {
			item.Channel = conf.Name
			if _, err := item.InsertIntoSQL(); err != nil {
				log.WithFields(log.Fields{
					"channel" : conf.Name,
					"page no." : i,
					"error" : err.Error(),
				}).Error("Errors when insert into sql")
				continue
			}
			cnt++
		}
	}

	log.WithFields(log.Fields{
		"channel" : conf.Name,
	}).Infof("End Catch, [%d] News expected, [%d] News Really Catched, Thread goes sleep", 20*conf.PagesOneTime, cnt)

}

func (conf *NeteaseNewsCatchConfigure)CatchDaemon() {
	for {
		select {
		case <-time.After(conf.DelayBetweenEachCatchRound) :
			conf.catchOneTime()
		case <-conf.quit :
			return
		}
	}
}

func (conf *NeteaseNewsCatchConfigure)Stop() {
	if conf.quit != nil {
		close(conf.quit)
	}
}

