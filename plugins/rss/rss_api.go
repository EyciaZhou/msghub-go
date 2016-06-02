package rss

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/EyciaZhou/msghub.go/interface"
	"github.com/EyciaZhou/msghub.go/plugins"
	"github.com/EyciaZhou/rss"
	"time"
)

type RssController struct {
	Url string `json:"url"`

	delayBetweenCatchRound time.Duration
}

func (p *RssController) GetFeed() (*rss.Feed, error) {
	feed, err := rss.Fetch(p.Url)
	if err != nil {
		return nil, err
	}
	return feed, err
}

type feed rss.Feed

func (p *feed) TransTo() (*Interface.Topic, error) {
	if p.Link == "" {
		return nil, errors.New("Link is empty")
	}
	result := new(Interface.Topic)

	_hash_link := md5.Sum(([]byte)(p.Link))

	result.Id = "rss_" + hex.EncodeToString(_hash_link[:])
	result.Title = p.Title

	msgs := make([]*Interface.Message, len(p.Items))

	author := &Interface.Author{}
	author.AvatarUrl = p.Image.Url
	author.Name = p.Nickname
	author.Uid = result.Id + "_" + author.Name

	lstModify := (int64)(0)

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
		next.Topic = result.Id
		next.Tag = ""
		next.Author = author
		//		next.Priority = 0

		if next.SnapTime > lstModify {
			lstModify = next.SnapTime
		}

		msgs[cnt] = next

		cnt++
	}
	result.LastModify = lstModify
	result.Msgs = msgs[:cnt]

	return result, nil
}

func (p *RssController) GetNew() (*Interface.Topic, error) {
	_feed, err := rss.Fetch(p.Url)
	if err != nil {
		return nil, err
	}
	p.delayBetweenCatchRound = _feed.Refresh.Sub(time.Now())
	return (*feed)(_feed).TransTo()
}

func (p *RssController) DelayBetweenCatchRound() time.Duration {
	return p.delayBetweenCatchRound
}

func NewRssController(url string) *RssController {
	return &RssController{
		url,
		time.Minute * 10,
	}
}

func LoadConf(conf_bs []byte) ([]plugins.GetNewer, error) {
	var confs []*RssController

	err := json.Unmarshal(conf_bs, &confs)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(confs); i++ {
		confs[i].delayBetweenCatchRound = time.Minute * 10
	}

	plugins := make([]plugins.GetNewer, len(confs))

	for i, conf := range confs {
		plugins[i] = conf
	}

	return plugins, nil
}

func init() {
	plugins.Register("Rss", (plugins.LoadConf)(LoadConf))
}
