package netease_news

import (
	"encoding/json"
	"errors"
	"github.com/EyciaZhou/msghub.go/ErrorUtiles"
	"github.com/EyciaZhou/msghub.go/plugins"
	"github.com/EyciaZhou/msghub.go/plugins/netease_news/api"
	"github.com/EyciaZhou/msghub.go/plugins/plugin"
	"time"
)

var (
	_DEFAULT_DELAY = 10 * time.Minute
)

type user_conf struct {
	Type string `can_null:"false" name:"类型" desc:"可为channel或special"`
	Id   string `can_null:"false" name:"ID" desc:"类型对应的id"`
}

func (p *PluginNeteaseNews) NewTask(config interface{}) (PluginInterface.PluginTask, error) {
	conf, ok := config.(*user_conf)
	if !ok {
		return nil, errors.New("config type error")
	}

	var result PluginInterface.PluginTask
	switch conf.Type {
	case "channel":
		if _, have := channelsDefault[conf.Id]; !have {
			return nil, ErrorUtiles.NewError("can't find builtin configure when type is defaultChannel")
			result = nenews_api.NewChannController(conf.Id, _DEFAULT_DELAY)
		}

		result = nenews_api.NewChannController(channelsDefault[conf.Id].ID, _DEFAULT_DELAY)
	case "special":
		result = nenews_api.NewTopicController(conf.Id, _DEFAULT_DELAY)
	default:
		return nil, errors.New("不支持的类型")
	}

	return result, nil
}

func (p *PluginNeteaseNews) GetConfigType() interface{} {
	return &user_conf{}
}

type resume_status struct {
	Type  string                 `json:"TYPE"`
	Value map[string]interface{} `json:"VALUE"`
}

func (p *PluginNeteaseNews) ResumeTask(Status []byte) (PluginInterface.PluginTask, error) {
	var rs resume_status
	json.Unmarshal(Status, &rs)
	if rs.Type == "SPECIAL" {
		return nenews_api.ResumeTopicController(rs.Value)
	} else if rs.Type == "CHANNEL" {
		return nenews_api.ResumeChannController(rs.Value)
	}
	return nil, errors.New("[neteasenews] : can't resume this pluginTask : unknow type : " + (string)(Status))
}

func (p *PluginNeteaseNews) Name() string {
	return "网易新闻"
}

type PluginNeteaseNews struct{}

var pluginNeteaseNews = &PluginNeteaseNews{}

func init() {
	plugin.Register("neteasenews", pluginNeteaseNews)
}
