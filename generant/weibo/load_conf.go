package weibo

import (
	"github.com/EyciaZhou/msghub.go/generant"
	"github.com/EyciaZhou/msghub.go/generant/weibo/api"
	"encoding/json"
	"time"
)

type conf struct {
	Token string `json:"token"`
	Delay string `json:"delay"`
	FetchEach int `json:"fetch_each"`
}

func LoadConf(raw []byte) ([]generant.GetNewer, error) {
	var c conf
	err := json.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}

	dur, err := time.ParseDuration(c.Delay)

	return []generant.GetNewer{
		weibo_api.NewFriendsTimelineController(c.Token, "", dur, c.FetchEach),
	}, nil
}

func init() {
	generant.Register("weibo", (generant.LoadConf)(LoadConf))
}
