package weibo_api

import (
	"github.com/EyciaZhou/msghub.go/generant/utils"
	"github.com/EyciaZhou/msghub.go/generant/weibo/types"
	"net/url"
	"github.com/EyciaZhou/msghub.go/generant"
	"time"
	"strconv"
)

const (
	_URL_FRIENDS_TIMELINE = "https://api.weibo.com/2/statuses/friends_timeline.json"
	_URL_ACCESS_TOKEN     = "https://api.weibo.com/oauth2/access_token"
	_URL_GET_TOKEN_INFO   = "https://api.weibo.com/oauth2/get_token_info"
	//_DEFAULT_WEIBO_DELAY = time.Minute * 10
	//_FETCH_EACH = 100
)

func (p *FriendsTimelineController) firstPage() ([]*weibo_types.Tweet, error) {
	args := url.Values{
		"access_token": {p.token},
		"count":        {strconv.Itoa(p.fetchEach)},
	}

	wtl := weibo_types.Timeline{}

	err := generant_utils.RequestToJson("GET", _URL_FRIENDS_TIMELINE, args, &wtl)

	if err != nil {
		return nil, err
	}

	if wtl.Statuses == nil {
		wtl.Statuses = []*weibo_types.Tweet{}
	}

	return wtl.Statuses, nil
}

func (p *FriendsTimelineController) since(SinceId string) ([]*weibo_types.Tweet, error) {
	tweets := []*weibo_types.Tweet{}
	for {
		tweets_new, err := p.pageFlip(SinceId)
		if err != nil {
			return nil, err
		}

		if len(tweets_new) == 0 {
			return tweets, nil
		}

		SinceId = tweets_new[0].Idstr

		tweets = append(tweets, tweets_new...)

		if len(tweets_new) < p.fetchEach / 2 {
			return tweets, nil
		}
	}
}

func (p *FriendsTimelineController) pageFlip(SinceId string) ([]*weibo_types.Tweet, error) {
	args := url.Values{
		"access_token": {p.token},
		"since_id":     {SinceId},
		"count":        {strconv.Itoa(p.fetchEach)},
	}

	wtl := weibo_types.Timeline{}

	err := generant_utils.RequestToJson("GET", _URL_FRIENDS_TIMELINE, args, &wtl)

	if err != nil {
		return nil, err
	}

	if wtl.Statuses == nil {
		wtl.Statuses = []*weibo_types.Tweet{}
	}

	return wtl.Statuses, nil
}

type FriendsTimelineController struct {
	token string
	lstid string

	delay time.Duration
	fetchEach int
}

func NewFriendsTimelineController(token string, lstid string, delay time.Duration, fetchEach int) *FriendsTimelineController {
	return &FriendsTimelineController{
		token: token,
		lstid: lstid,
		delay: delay,
		fetchEach: fetchEach,
	}
}

func (p *FriendsTimelineController) DelayBetweenCatchRound() time.Duration {
	return p.delay
}

func (p *FriendsTimelineController) GetNew() (generant.CanConvertToTopic, error) {
	if p.lstid == "" {
		ts, err := p.firstPage()
		if err == nil && len(ts) > 0 {
			p.lstid = ts[0].Idstr
		}
		return (weibo_types.Tweets)(ts), err
	}

	ts, err := p.since(p.lstid)
	if err == nil && len(ts) > 0 {
		p.lstid = ts[0].Idstr
	}
	return (weibo_types.Tweets)(ts), err
}
