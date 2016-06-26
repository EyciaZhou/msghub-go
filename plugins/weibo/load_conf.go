package weibo

import (
	"errors"
	"github.com/EyciaZhou/msghub.go/plugins"
	"github.com/EyciaZhou/msghub.go/plugins/plugin"
	"github.com/EyciaZhou/msghub.go/plugins/weibo/api"
	"time"
)

type user_conf struct {
	Token string `can_null:"false" desc:"微博的token"`
}

func (p *PluginWeibo) ResumeTask(status []byte) (PluginInterface.PluginTask, error) {
	return weibo_api.ResumeFriendsTimelineController(status)
}

func (p *PluginWeibo) NewTask(Config interface{}) (PluginInterface.PluginTask, error) {
	if v, ok := Config.(*user_conf); !ok {
		return nil, errors.New("config type error")
	} else {
		return weibo_api.NewFriendsTimelineController(
			v.Token,
			"",
			10*time.Minute,
			100,
		), nil
	}
}

func (p *PluginWeibo) GetConfigType() interface{} {
	return &user_conf{}
}

func (p *PluginWeibo) Name() string {
	return "新浪微博"
}

type PluginWeibo struct{}

var pluginWeibo = &PluginWeibo{}

func init() {
	plugin.Register("weibo", pluginWeibo)
}
