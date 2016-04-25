package nenews_types

import (
	"errors"
	"github.com/EyciaZhou/msghub.go/generant"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"time"
)

var (
	//2015-09-02 08:33:27
	timeFormat = "2006-01-02 15:04:05"
	loc        *time.Location
)

var (
	authorNetEaseNews = &generant.Author{
		Uid:         "NetEaseNews",
		Name:        "网易新闻",
		CoverSource: "http://www.apk20.com/image/icon-385074",
	}
)

func parseTime(ts string) (int64, error) {
	ti, e := time.ParseInLocation(timeFormat, ts, loc)
	if e != nil {
		return 0, e
	}
	return ti.Unix(), nil
}

func (r Reply) Convert() (generant.Reply, error) {
	var e error

	length := len(r)
	reply := make(generant.Reply, length)

	cnt := 0
	for i := 1; i <= length; i++ {
		id := strconv.Itoa(i)
		if _, ok := r[id]; !ok {
			return nil, errors.New("error when parse reply, missing floor " + id)
		}

		ti := (int64)(0)
		dig := 0

		if i == length {
			//only last floor have time and digg
			ti, e = parseTime(r[id].Time)
			if e != nil {
				return nil, e
			}
			dig, e = strconv.Atoi((r)[id].Digg)
			if e != nil {
				return nil, e
			}
		}

		reply[cnt] = &generant.ReplyFloor{
			Floor:   i,
			Time:    ti,
			Name:    r[id].Name,
			Content: r[id].Content,
			Digg:    dig,
		}
		cnt++
	}
	return reply, nil
}

func (p *NewsImage) Convert() *generant.Image {
	return &generant.Image{
		p.Ref, p.Title, p.Size, p.URL,
	}
}

func (p *PhototSetImage) Convert() *generant.Image {
	return &generant.Image{
		"", p.Desc, "", p.URL,
	}
}

func (n *PhotoSet) Convert() (*generant.Message, error) {
	var replys []generant.Reply
	var imgs []*generant.Image

	//process reply
	for _, item := range n.Replys {
		nReply, err := item.Convert()
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		replys = append(replys, nReply)
	}

	//process images
	for _, item := range n.Images {
		nImage := item.Convert()
		imgs = append(imgs, nImage)
	}
	pubti, err := parseTime(n.PubTime)
	if err != nil {
		return nil, err
	}

	snapti, err := parseTime(n.SnapTime)
	if err != nil {
		snapti = time.Now().Unix()
	}

	return &generant.Message{
		SnapTime:    snapti,
		PubTime:     pubti,
		Source:      n.URL,
		Body:        n.Body,
		Title:       n.Title,
		Subtitle:    n.SubTitle,
		CoverImg:    n.CoverURL,
		Images:      imgs,
		ReplyNumber: (int64)(n.ReplyCount),
		Replys:      replys,
		ViewType:    generant.VIEW_TYPE_PICTURES,
		Author:      authorNetEaseNews,
	}, nil
}

func (n *NormalNews) Convert() (*generant.Message, error) {
	var replys []generant.Reply
	var imgs []*generant.Image

	//progress reply
	for _, item := range n.Replys {
		nReply, err := item.Convert()
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		replys = append(replys, nReply)
	}

	//progress imgages
	for _, item := range n.Images {
		nImage := item.Convert()

		imgs = append(imgs, nImage)
	}

	pubti, err := parseTime(n.PubTime)
	if err != nil {
		return nil, err
	}

	snapti, err := parseTime(n.SnapTime)
	if err != nil {
		snapti = time.Now().Unix()
	}

	return &generant.Message{
		SnapTime:    snapti,
		PubTime:     pubti,
		Source:      n.URL,
		Body:        n.Body,
		Title:       n.Title,
		Subtitle:    n.SubTitle,
		CoverImg:    n.CoverURL,
		Images:      imgs,
		ReplyNumber: (int64)(n.ReplyCount),
		Replys:      replys,
		ViewType:    generant.VIEW_TYPE_NORMAL,
		Version:     "0.1",
		Author:      authorNetEaseNews,
	}, nil
}

func (p *Topic)Convert() *generant.Topic {
	result := make([]*generant.Message, len(p.Newss))

	cnt := 0

	for _, news := range p.Newss {
		m, e := news.Convert()
		if e != nil {
			log.Warn(e.Error())
			continue
		}
		result[cnt] = m
		cnt++
	}

	return &generant.Topic{
		Id:p.Id,
		Title:p.Title,
		Msgs:result[:cnt],
		LastModify:time.Now().Unix(),
	}
}

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Error(err.Error())
	}
}
