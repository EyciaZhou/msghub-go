package weibo_api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/EyciaZhou/msghub-http/Utils"
	"github.com/EyciaZhou/msghub.go/HttpUtils"
	"github.com/EyciaZhou/msghub.go/interface"
	"github.com/EyciaZhou/msghub.go/plugins/weibo/types"
	"net/url"
	"strconv"
	"time"
)

const (
	_URL_FRIENDS_TIMELINE = "https://api.weibo.com/2/statuses/friends_timeline.json"
	_URL_ACCESS_TOKEN     = "https://api.weibo.com/oauth2/access_token"
	_URL_GET_TOKEN_INFO   = "https://api.weibo.com/oauth2/get_token_info"

//_DEFAULT_WEIBO_DELAY = time.Minute * 10
//_FETCH_EACH = 100
)

func weiboErrorChecker(bs []byte) error {
	we := weibo_types.Error{}
	err1 := json.Unmarshal(bs, &we)
	if err1 != nil {
		return err1
	}
	if we.Error != "" {
		return errors.New(we.Error)
	}
	return nil
}

func (p *FriendsTimelineController) firstPage() ([]*weibo_types.Tweet, error) {
	args := url.Values{
		"access_token": {p.Token},
		"count":        {strconv.Itoa(p.FetchEach)},
	}

	wtl := weibo_types.Timeline{}

	err := HttpUtils.JsonCheckError("GET", _URL_FRIENDS_TIMELINE, args, weiboErrorChecker, &wtl)

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

		if len(tweets_new) < p.FetchEach/2 {
			return tweets, nil
		}
	}
}

func (p *FriendsTimelineController) pageFlip(SinceId string) ([]*weibo_types.Tweet, error) {
	args := url.Values{
		"access_token": {p.Token},
		"since_id":     {SinceId},
		"count":        {strconv.Itoa(p.FetchEach)},
	}

	wtl := weibo_types.Timeline{}

	err := HttpUtils.JsonCheckError("GET", _URL_FRIENDS_TIMELINE, args, weiboErrorChecker, &wtl)

	if err != nil {
		return nil, err
	}

	if wtl.Statuses == nil {
		wtl.Statuses = []*weibo_types.Tweet{}
	}

	return wtl.Statuses, nil
}

type FriendsTimelineController struct {
	Token string `json:"token"`
	Lstid string `json:"lstid"`

	Delay     time.Duration `json:"delay"`
	FetchEach int           `json:"fetch_each"`
}

func ResumeFriendsTimelineController(Status []byte) (*FriendsTimelineController, error) {
	result := &FriendsTimelineController{}
	err := json.Unmarshal(Status, result)
	return result, err
}

func NewFriendsTimelineController(token string, lstid string, delay time.Duration, fetchEach int) *FriendsTimelineController {
	return &FriendsTimelineController{
		Token:     token,
		Lstid:     lstid,
		Delay:     delay,
		FetchEach: fetchEach,
	}
}

func (p *FriendsTimelineController) Type() string {
	return "weibo"
}

func (p *FriendsTimelineController) GetDelayBetweenCatchRound() time.Duration {
	return p.Delay
}

func (p *FriendsTimelineController) FetchNew() ([]*Interface.Message, error) {
	if p.Lstid == "" {
		ts, err := p.firstPage()
		if err == nil && len(ts) > 0 {
			p.Lstid = ts[0].Idstr
		}
		return (weibo_types.Tweets)(ts).Convert(), err
	}

	ts, err := p.since(p.Lstid)
	if err == nil && len(ts) > 0 {
		p.Lstid = ts[0].Idstr
	}
	return (weibo_types.Tweets)(ts).Convert(), err
}

func (p *FriendsTimelineController) DumpTaskStatus() (Status []byte) {
	result, _ := json.Marshal(p)
	return result
}

func (p *FriendsTimelineController) Hash() string {
	return "weibo_" + hex.EncodeToString(Utils.Sha1(([]byte)(p.Token)))
}
