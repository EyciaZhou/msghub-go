package weibo

import (
	"encoding/json"
	"github.com/EyciaZhou/msghub.go/plugins"
	"github.com/EyciaZhou/msghub.go/plugins/weibo/api"
	"time"
)

type conf struct {
	Token     string `json:"token"`
	Delay     string `json:"delay"`
	FetchEach int    `json:"fetch_each"`
}

func LoadConf(raw []byte) ([]plugins.GetNewer, error) {
	var c conf
	err := json.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}

	dur, err := time.ParseDuration(c.Delay)

	return []plugins.GetNewer{
		weibo_api.NewFriendsTimelineController(c.Token, "", dur, c.FetchEach),
	}, nil
}

func init() {
	plugins.Register("weibo", (plugins.LoadConf)(LoadConf))
}
