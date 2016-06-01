package netease_news

import (
	"encoding/json"
	"github.com/EyciaZhou/msghub.go/ErrorUtiles"
	"github.com/EyciaZhou/msghub.go/plugins"
	"github.com/EyciaZhou/msghub.go/plugins/netease_news/api"
	"time"
)

var (
	_DEFAULT_DELAY = 10 * time.Minute
)

func LoadConf(conf_bs []byte) ([]plugins.GetNewer, error) {
	confs := []map[string]string{}
	err := json.Unmarshal(conf_bs, &confs)
	if err != nil {
		return nil, ErrorUtiles.NewPanicError(err)
	}

	result := []plugins.GetNewer{}
	for _, conf := range confs {
		delayTime_S := conf["delayBetweenCatchRound"]

		delayTime, err := time.ParseDuration(delayTime_S)
		if err != nil {
			delayTime = _DEFAULT_DELAY
		}

		var (
			gn plugins.GetNewer
		)

		switch conf["type"] {
		case "defaultChannel":
			channelName := conf["channelName"]
			if _, have := channelsDefault[channelName]; !have {
				return nil, ErrorUtiles.NewError("can't find builtin configure when type is defaultChannel")
			}
			gn = nenews_api.NewChannController(channelsDefault[channelName].Name, channelsDefault[channelName].ID, delayTime)
		case "customChannel":
			name := conf["name"]
			id := conf["id"]
			if name == "" || id == "" {
				return nil, ErrorUtiles.NewError("name or id is empty when type is customChannel")
			}
			gn = nenews_api.NewChannController(name, id, delayTime)
		case "special":
			id := conf["topicId"]
			if id == "" {
				return nil, ErrorUtiles.NewError("topicId is empty when type is special")
			}
			gn = nenews_api.NewTopicController(id, delayTime)
		default:
			return nil, ErrorUtiles.NewError("not supported type of nenews conf")
		}

		result = append(result, gn)
	}
	return result, nil
}

func init() {
	plugins.Register("NeteaseNews", (plugins.LoadConf)(LoadConf))
}
