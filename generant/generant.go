package generant

import (
	"github.com/Sirupsen/logrus"
	"time"
)

type Generant struct {
	GetNewer

	forceStop chan struct{}
	exited chan struct{}
}

func NewGrenrant(getNewer GetNewer) *Generant {
	return &Generant{
		getNewer,
		make(chan struct{}, 1),
		make(chan struct{}, 1),
	}
}

func (p *Generant) catchOneTime(roundEnd chan bool) {
	logrus.Info("start catch")

	FetchIsSucc := make(chan bool, 1)
	var (
		m CanConvertToTopic
		e error
	)

	go func() {
		m, e = p.GetNewer.GetNew()
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

	topic := m.Convert()

	e = topic.InsertIntoSQL()
	if e != nil {
		logrus.Errorf("Insert Topic[%s] into sql Error: %s", topic.Id, e.Error())
	}
	roundEnd <- true
	logrus.Infof("%d catched", len(topic.Msgs))
}

func (p *Generant) Catch() {
	go func() {
		defer func() {close(p.exited)}()

		oneCatchEnded := make(chan bool, 1)
		for {
			go p.catchOneTime(oneCatchEnded)

			select {
			case <-oneCatchEnded:
			case <-p.forceStop:
				return
			}

			select {
			case <-time.After(60 * time.Second):
			case <-p.forceStop:
				return
			}
		}
	}()
}

func (p *Generant) Stop() chan struct{} {
	close(p.forceStop)
	return p.exited
}