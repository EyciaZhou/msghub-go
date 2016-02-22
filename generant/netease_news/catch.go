package netease_news

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"errors"
	"git.eycia.me/eycia/msghub/generant"
)

type NeteaseNewsConfigure interface {
	DelayBetweenEachCatchRound() time.Duration
	CatchOneTime()
	Quit() chan int
}

type NeteaseNewsTopicConfigure struct {
	TopicId string
	delayBetweenEachCatchRound time.Duration

	quit chan int
}

func NewNeteaseNewsTopicConfigure(topicid string, delayBetweenEachCatchRound time.Duration) (*NeteaseNewsTopicConfigure, error) {
	return &NeteaseNewsTopicConfigure{
		TopicId:topicid,
		delayBetweenEachCatchRound:delayBetweenEachCatchRound,
		quit:make(chan int),
	}, nil
}

type NeteaseNewsChannelConfigure struct {
	NeteaseNewsChannel
	pagesOneTime		int
	delayBetweenPage	time.Duration
	delayBetweenEachCatchRound time.Duration

	quit chan int
}

func NewNeteaseNewsCatchConfigure(chann NeteaseNewsChannel, pagesOneTime int,
delayBetweenPage, delayBetweenEachCatchRound time.Duration) (*NeteaseNewsChannelConfigure, error) {
	return &NeteaseNewsChannelConfigure{
		chann,
		pagesOneTime,
		delayBetweenPage,
		delayBetweenEachCatchRound,

		make(chan int),
	}, nil
}

func NewDefaultNeteaseNewsChannelConfigure(channelName string) (*NeteaseNewsChannelConfigure, error) {
	if _, hv := channelsDefault[channelName]; !hv {
		return nil, errors.New("no such channel")
	}
	channelInfo := channelsDefault[channelName]
	return &NeteaseNewsChannelConfigure{
		*channelInfo,
		2,
		time.Second * 10,
		time.Minute * 10,

		make(chan int),
	}, nil
}

func (conf *NeteaseNewsTopicConfigure)CatchOneTime() {
	chanFetch := make(chan int, 1)
	var (
		m *generant.Topic
		e error
	)

	go func() {
		m, e = getSpecialList(conf.TopicId)
		if e != nil {
			log.Errorf("Fetch Topic[%s] Error: %s", conf.TopicId, e.Error())
			chanFetch <- -1
			return
		}
		chanFetch <- 1
	}()

	//block waiting fetch complete or quit command, when receive quit command during fetching,
	//leave it fetching, and this function return, when receive quit command during Insert,
	//wait Insert complete
	select {
	case res := <-chanFetch:
		if res < 0 {
			return
		}

	case <-conf.quit:
		return
	}

	e = m.InsertIntoSQL()
	if e != nil {
		log.Errorf("Insert Topic[%s] into sql Error: %s", conf.TopicId, e.Error())
	}
	log.Infof("Topic[%s] fetch finished", conf.TopicId)
}

func (conf *NeteaseNewsTopicConfigure)Quit() chan int {
	return conf.quit
}

func (conf *NeteaseNewsTopicConfigure)DelayBetweenEachCatchRound() time.Duration {
	return conf.delayBetweenEachCatchRound
}

func (conf *NeteaseNewsChannelConfigure)Quit() chan int {
	return conf.quit
}

func (conf *NeteaseNewsChannelConfigure)DelayBetweenEachCatchRound() time.Duration {
	return conf.delayBetweenEachCatchRound
}

func (conf *NeteaseNewsChannelConfigure)CatchOneTime() {
	log.WithFields(log.Fields{
		"channel" : conf.Name,
		"page num" : conf.pagesOneTime,
	}).Info("Start Catch")

	cnt := 0

	for i := 0; i < conf.pagesOneTime; i++ {
		if i != 0 {
			select {
			case <-time.After(conf.delayBetweenPage) :

			case <-conf.quit :
				return
			}
		}

		log.WithFields(log.Fields{
			"channel" : conf.Name,
			"page no." : i,
		}).Info("Start Catch Page")

		news, err := getNewsChannel(conf.ID, i)
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
	}).Infof("End Catch, [%d] News expected, [%d] News Really Catched, Thread goes sleep", 20*conf.pagesOneTime, cnt)

}

func catchDaemon(conf NeteaseNewsConfigure) {
	conf.CatchOneTime()
	for {
		select {
		case <-time.After(conf.DelayBetweenEachCatchRound()) :
			conf.CatchOneTime()
		case <-conf.Quit():
			return
		}
	}
}

func stopCatchDaemon(conf NeteaseNewsConfigure) {
	if conf.Quit() != nil {
		close(conf.Quit())
	}
}

