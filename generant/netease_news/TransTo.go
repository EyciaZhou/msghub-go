package netease_news

import (
	log "github.com/Sirupsen/logrus"
	"git.eycia.me/eycia/msghub/generant"
	"strconv"
	"time"
)

var (
	//2015-09-02 08:33:27
	timeFormat = "2006-01-02 15:04:05"
	loc        *time.Location
)

func parseTime(ts string) (int64, error) {
	ti, e := time.ParseInLocation(timeFormat, ts, loc)
	if e != nil {
		return 0, e
	}
	return ti.Unix(), nil
}

//ToReply: if any floor errors, it returns nil and error
func (r Reply) ToReply() (generant.Reply, error) {
	var reply generant.Reply

	length := len(r)
	for i := 1; i <= length; i++ {
		id := strconv.Itoa(i)
		ti, _ := parseTime(r[id].Time)
		/* only last floor have time
		if e != nil {
			return nil, e
		}
		*/
		dig, _ := strconv.Atoi((r)[id].Digg)
		/*
			only last florr have digg
			if e != nil {
				return nil, e
			}
		*/
		reply = append(reply, &generant.ReplyFloor{
			Floor:   i,
			Time:    ti,
			Name:    (r)[id].Name,
			Content: (r)[id].Content,
			Digg:    dig,
		})
	}
	return reply, nil
}

func (p *NewsImage) ToImage() (*generant.Image, error) {
	return &generant.Image{
		p.Ref, p.Title, p.URL,
	}, nil
}

func (p *PhototSetImage) ToImage() (*generant.Image, error) {
	return &generant.Image{
		"", p.Desc, p.URL,
	}, nil
}

func (n *PhotoSet) ToMsg() (*generant.Message, error) {
	var replys []generant.Reply
	var imgs []*generant.Image

	//process reply
	for _, item := range n.Replys {
		nReply, err := item.ToReply()
		if err != nil {
			log.WithFields(log.Fields{
				"time" : "fetch",
				"reply" : item,
				"error" : err.Error(),
			}).Warning("error throwed when reply trans")
			continue
		}
		replys = append(replys, nReply)
	}

	//process images
	for _, item := range n.Images {
		nImage, err := item.ToImage()
		if err != nil {
			log.WithFields(log.Fields{
				"time" : "fetch",
				"image" : item,
				"error" : err.Error(),
			}).Warning("error throwed when image trans")
			continue
		}
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
		ReplyNumber: n.ReplyCount,
		Replys:      replys,
		ViewType:    n.ViewType,
		Version:     "0.1",
		From:        "Netease News",
		Priority:    n.Priority,
	}, nil
}

func (n *News) ToMsg() (*generant.Message, error) {
	var replys []generant.Reply
	var imgs []*generant.Image

	//progress reply
	for _, item := range n.Replys {
		nReply, err := item.ToReply()
		if err != nil {
			log.WithFields(log.Fields{
				"time" : "fetch",
				"reply" : item,
				"error" : err.Error(),
			}).Warning("error throwed when reply trans")
			continue
		}
		replys = append(replys, nReply)
	}

	//progress imgages
	for _, item := range n.Images {
		nImage, err := item.ToImage()
		if err != nil {
			log.WithFields(log.Fields{
				"time" : "fetch",
				"image" : item,
				"error" : err.Error(),
			}).Warning("error throwed when image trans")
			continue
		}

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
		ReplyNumber: n.ReplyCount,
		Replys:      replys,
		ViewType:    n.ViewType,
		Version:     "0.1",
		From:        "Netease News",
		Priority:    n.Priority,
	}, nil
}

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Error(err.Error())
	}
}
