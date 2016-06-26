package rss

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/EyciaZhou/msghub.go/interface"
	"github.com/EyciaZhou/msghub.go/plugins"
	"github.com/EyciaZhou/msghub.go/plugins/plugin"
	"github.com/EyciaZhou/rss"
	"time"
)

type user_conf struct {
	URL string `can_null:"false" desc:"RSS地址"`
}

type RssController struct {
	Url string `json:"url"`

	delayBetweenCatchRound time.Duration `json:"delay"`
}

func (p *RssController) Type() string {
	return "rss"
}

func (p *RssController) GetFeed() (*rss.Feed, error) {
	feed, err := rss.Fetch(p.Url)
	if err != nil {
		return nil, err
	}
	return feed, err
}

type feed rss.Feed

func (p *feed) TransTo(FecherId string) ([]*Interface.Message, error) {
	if p.Link == "" {
		return nil, errors.New("Link is empty")
	}

	msgs := make([]*Interface.Message, len(p.Items))

	author := &Interface.Author{}
	author.AvatarUrl = p.Image.Url
	author.Name = p.Nickname
	author.Uid = FecherId + "_" + author.Name

	cnt := 0
	for _, item := range p.Items {
		if item.Link == "" {
			continue
		}

		next := new(Interface.Message)
		next.SnapTime = item.Date.Unix()
		next.PubTime = item.Date.Unix()
		next.Source = item.Link
		next.Body = item.Content
		next.Title = item.Title
		next.Subtitle = item.Summary
		next.CoverImg = "" //?
		next.Images = nil
		next.ViewType = Interface.VIEW_TYPE_NORMAL
		next.Tag = ""
		next.Author = author
		//		next.Priority = 0

		msgs[cnt] = next

		cnt++
	}
	return msgs[:cnt], nil
}

func (p *RssController) FetchNew() ([]*Interface.Message, error) {
	_feed, err := rss.Fetch(p.Url)
	if err != nil {
		return nil, err
	}
	p.delayBetweenCatchRound = _feed.Refresh.Sub(time.Now())
	return (*feed)(_feed).TransTo(p.Hash())
}

func (p *RssController) GetDelayBetweenCatchRound() time.Duration {
	return p.delayBetweenCatchRound
}

func (p *RssController) DumpTaskStatus() (Status []byte) {
	bs, _ := json.Marshal(p)
	return bs
}

func (p *RssController) Hash() string {
	md5ed := md5.Sum(([]byte(p.Url)))
	return "rss_" + hex.EncodeToString(md5ed[:])
}

func NewRssController(url string) *RssController {
	return &RssController{
		url,
		time.Minute * 10,
	}
}

type PluginRss struct{}

var pluginRss = &PluginRss{}

func (p *PluginRss) ResumeTask(status []byte) (PluginInterface.PluginTask, error) {
	confs := RssController{}
	err := json.Unmarshal(status, &confs)
	if err != nil {
		return nil, err
	}
	return &confs, nil
}

func (p *PluginRss) NewTask(Config interface{}) (PluginInterface.PluginTask, error) {
	if v, ok := Config.(*user_conf); !ok {
		return nil, errors.New("config type error")
	} else {
		return NewRssController(v.URL), nil
	}
}

func (p *PluginRss) GetConfigType() interface{} {
	return &user_conf{}
}

func (p *PluginRss) Name() string {
	return "RSS订阅"
}

func init() {
	plugin.Register("rss", pluginRss)
}
