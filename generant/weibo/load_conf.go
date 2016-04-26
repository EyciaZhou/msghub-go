package weibo

import (
	"github.com/EyciaZhou/msghub.go/generant"
	"github.com/EyciaZhou/msghub.go/generant/weibo/api"
)

func LoadConf(raw []byte) ([]generant.GetNewer, error) {
	return []generant.GetNewer{
		weibo_api.NewFriendsTimelineController((string)(raw), ""),
	}, nil
}

func init() {
	generant.Register("NeteaseNews", (generant.LoadConf)(LoadConf))
}
