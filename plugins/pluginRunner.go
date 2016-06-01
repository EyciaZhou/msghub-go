package plugins

import (
	"github.com/EyciaZhou/msghub.go/interface"
	"github.com/Sirupsen/logrus"
	"time"
)

type pluginRunner struct {
	GetNewer

	forceStop chan struct{}
	exited    chan struct{}
}

func NewPluginRunner(getNewer GetNewer) *pluginRunner {
	return &pluginRunner{
		getNewer,
		make(chan struct{}, 1),
		make(chan struct{}, 1),
	}
}

func (p *pluginRunner) catchOneTime(roundEnd chan struct{}) {
	logrus.Info("start catch")
	defer func() { roundEnd <- struct{}{} }()

	FetchIsSucc := make(chan bool, 1)
	var (
		topic *Interface.Topic
		e     error
	)

	go func() {
		topic, e = p.GetNewer.GetNew()
		if e != nil {
			logrus.Error(e)
			FetchIsSucc <- false
			return
		}
		FetchIsSucc <- true
	}()

	//block waiting fetch complete or quit command, when receive quit command during fetching,
	//leave it fetching, and this function return, when receive quit command during Insert,
	//wait Insert complete
	select {
	case succ := <-FetchIsSucc:
		if !succ {
			return
		}

	case <-p.forceStop:
		return
	}

	e = topic.InsertIntoSQL()
	if e != nil {
		logrus.Errorf("Insert Topic[%s] into sql Error: %s", topic.Id, e.Error())
	}
	logrus.Infof("%d catched", len(topic.Msgs))
}

func (p *pluginRunner) Catch() {
	go func() {
		defer func() { close(p.exited) }()

		oneCatchEnded := make(chan struct{}, 1)
		for {
			go p.catchOneTime(oneCatchEnded)

			select {
			case <-oneCatchEnded:
			case <-p.forceStop:
				return
			}

			select {
			case <-time.After(p.GetNewer.DelayBetweenCatchRound()):
			case <-p.forceStop:
				return
			}
		}
	}()
}

func (p *pluginRunner) Stop() chan struct{} {
	close(p.forceStop)
	return p.exited
}
