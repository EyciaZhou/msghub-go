package weibo

import (
	"github.com/EyciaZhou/msghub.go/generant"
	"github.com/EyciaZhou/msghub.go/generant/weibo/generant"
)

func LoadConf(raw []byte) ([]generant.Generant, error) {
	//pass now

	return []generant.Generant{
		weibo_generant.NewFriendsTimelineGrenrant(),
	}, nil
}

func init() {
	generant.Register("weibo", (generant.LoadConf)(LoadConf))
}
