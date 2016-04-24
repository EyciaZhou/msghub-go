package weibo_api

import (
	"github.com/EyciaZhou/msghub.go/generant/weibo/types"
	"github.com/EyciaZhou/msghub.go/generant/weibo/utils"
	"net/url"
)

const (
	_URL_FRIENDS_TIMELINE = "https://api.weibo.com/2/statuses/friends_timeline.json"
	_URL_ACCESS_TOKEN     = "https://api.weibo.com/oauth2/access_token"
	_URL_GET_TOKEN_INFO   = "https://api.weibo.com/oauth2/get_token_info"
)

func (p *FriendsTimelineController)firstPage() ([]*weibo_types.Tweet, error) {
	args := url.Values{
		"access_token": {p.token},
		"count":        {"100"},
	}

	wtl := weibo_types.Timeline{}

	err := weibo_utils.RequestToJson("GET", _URL_FRIENDS_TIMELINE, args, &wtl)

	if err != nil {
		return nil, err
	}

	if wtl.Statuses == nil {
		wtl.Statuses = []*weibo_types.Tweet{}
	}

	return wtl.Statuses, nil
}

func (p *FriendsTimelineController)since(SinceId string) ([]*weibo_types.Tweet, error) {
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
	}
}

func (p *FriendsTimelineController)pageFlip(SinceId string) ([]*weibo_types.Tweet, error) {
	args := url.Values{
		"access_token": {p.token},
		"since_id":     {SinceId},
		"count":        {"100"},
	}

	wtl := weibo_types.Timeline{}

	err := weibo_utils.RequestToJson("GET", _URL_FRIENDS_TIMELINE, args, &wtl)

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
}

func NewFriendsTimelineController(token string, lstid string) *FriendsTimelineController {
	return &FriendsTimelineController{
		token : token,
		lstid : lstid,
	}
}

func (p *FriendsTimelineController) GetNew()  (weibo_types.Tweets, error) {
	if p.lstid == "" {
		ts, err := p.firstPage()
		if err == nil && len(ts) > 0 {
			p.lstid = ts[0].Idstr
		}
		return ts, err
	}

	ts, err := p.since(p.lstid)
	if err == nil && len(ts) > 0 {
		p.lstid = ts[0].Idstr
	}
	return ts, err
}
