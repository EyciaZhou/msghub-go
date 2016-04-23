package weibo_generant
import (
	"github.com/EyciaZhou/msghub.go/generant/weibo/api"
	"time"
	"github.com/EyciaZhou/msghub.go/generant/weibo/types"
	"github.com/Sirupsen/logrus"
)

var (
	_TOKEN       = "2.00KL_4zDiyz9HDe756a60853a_iJ3D" //valid in 5 years since 16y/4m/18d
)

type friendsTimelineGenerant struct {
	friendsline *weibo_api.FriendsTimelineController

	forceStop chan bool
}

func NewFriendsTimelineGrenrant() *friendsTimelineGenerant {
	return &friendsTimelineGenerant{
		weibo_api.NewFriendsTimelineController(_TOKEN, ""),
		make(chan bool, 1),
	}
}

func (p *friendsTimelineGenerant) catchOneTime(ended chan bool) {
	chanFetch := make(chan int, 1)
	var (
		m weibo_types.Tweets
		e error
	)

	go func() {
		m, e = p.friendsline.GetNew()
		if e != nil {
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

	case <-p.forceStop:
		return
	}

	topic := m.Convert()

	e = topic.InsertIntoSQL()
	if e != nil {
		logrus.Errorf("Insert Topic[%s] into sql Error: %s", topic.Id, e.Error())
	}
	ended <- true
}

func (p *friendsTimelineGenerant) Catch() {
	go func() {
		ended := make(chan bool, 1)
		for ;; {
			go p.catchOneTime(ended)
			select {
			case <-ended:
			case <-p.forceStop:
				return
			}

			select {
			case <-time.After(20 * time.Second):
			case <-p.forceStop:
				return
			}
		}
	}()
}

func (p *friendsTimelineGenerant) Stop() {
	close(p.forceStop)
}

func (p *friendsTimelineGenerant) ForceStop() {
	p.Stop()
}