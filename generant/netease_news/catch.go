package netease_news

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

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
	conf.catchOneTime()
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

