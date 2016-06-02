package weibo_types

import (
	"github.com/EyciaZhou/msghub.go/interface"
	"github.com/Sirupsen/logrus"
	"strings"
	"time"
)

type User struct {
	Name       string `json:"name"`
	Idstr      string `json:"idstr"`
	CoverImage string `json:"avatar_large"`
}

type PicUrl struct {
	ThumbnailPic string `json:"thumbnail_pic"`
}

type PicUrls []PicUrl

type Tweet struct {
	CreatedAt       string  `json:"created_at"`
	Id              int64   `json:"id"`
	Mid             string  `json:"mid"`
	Idstr           string  `json:"idstr"`
	Text            string  `json:"text"`
	PicUrls         PicUrls `json:"pic_urls"`
	ThumbnailPic    string  `json:"thumbnail_pic"`
	RepostsCount    int64   `json:"reposts_count"`
	CommentsCount   int64   `json:"comments_count"`
	AttitudesCount  int64   `json:"attitudes_count"`
	RetweetedStatus *Tweet  `json:"retweeted_status"`
	User            *User   `json:"user"`
}

type Timeline struct {
	Statuses []*Tweet `json:"statuses"`
}

type Tweets []*Tweet

func (p Tweets) Convert() *Interface.Topic {
	result := make([]*Interface.Message, len(p))

	cnt := 0

	for _, tweet := range p {
		m, e := tweet.Convert()
		if e != nil {
			logrus.Warn(e.Error())
			continue
		}
		result[cnt] = m
		cnt++
	}

	return &Interface.Topic{
		"weibo_friendsline",
		"weibo",
		result[:cnt],
		time.Now().Unix(),
	}
}

type Error struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Request   string `json:"request"`
}

var (
	weiboTimeLayout = time.RubyDate
)

const (
	_BASE62_KEYBOARD = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func (p *Tweet) getMidBase62() string {
	mid := p.Id

	result := ""

	for mid > 0 {
		k := mid % 10000000
		mid /= 10000000

		for k > 0 {
			result = string(_BASE62_KEYBOARD[k%62]) + result
			k /= 62
		}
	}

	return result
}

func (p *Tweet) GetSource() string {
	return "http://weibo.com/" + p.User.Idstr + "/" + p.getMidBase62()
}

func (p PicUrls) Convert() []*Interface.Image {
	if p == nil {
		return nil
	}

	result := make([]*Interface.Image, len(p))
	for index, img_from := range p {
		result[index] = &Interface.Image{
			URL: strings.Replace(img_from.ThumbnailPic, "thumbnail", "large", 1), //covert to big picture
		}
	}

	return result
}

func (p *User) Convert() *Interface.Author {
	return &Interface.Author{
		Name:      p.Name,
		Uid:       "weibo_" + p.Idstr,
		AvatarUrl: p.CoverImage,
	}
}

func (p *Tweet) Convert() (*Interface.Message, error) {
	result := &Interface.Message{}

	pubtime, err := time.Parse(weiboTimeLayout, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	result.SnapTime = pubtime.Unix()
	result.PubTime = pubtime.Unix()
	result.Source = p.GetSource()
	result.Body = p.Text

	if p.RetweetedStatus != nil {
		result.Body += "//"

		if p.RetweetedStatus.User != nil {
			result.Body += "@" + p.RetweetedStatus.User.Name + ":"
		}

		result.Body += p.RetweetedStatus.Text
	}

	result.Title = ""             //no title, replace with author's name
	result.Subtitle = result.Body //display on first screen
	//	result.ReplyNumber = p.CommentsCount
	//	result.Replys = nil //TODO
	result.ViewType = Interface.VIEW_TYPE_PICTURES
	result.Topic = "weibo_friendsline"
	result.Author = p.User.Convert()
	result.Tag = "" //TODO

	//if retweet, use Retweeted's coverimg and imgages
	if p.RetweetedStatus != nil {
		p = p.RetweetedStatus
	}

	result.CoverImg = p.ThumbnailPic //"" if not have
	result.Images = p.PicUrls.Convert()

	return result, nil
}
