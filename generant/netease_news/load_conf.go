package netease_news
import (
	"github.com/EyciaZhou/msghub.go/generant"
	"encoding/json"
	"github.com/EyciaZhou/msghub.go/Utiles"
	"time"
	"github.com/EyciaZhou/msghub.go/generant/netease_news/api"
	"github.com/Sirupsen/logrus"
)

var (
	_DEFAULT_DELAY = 10 * time.Minute
)

func LoadConf(conf_bs []byte) ([]generant.GetNewer, error) {
	confs := []map[string]string{}
	err := json.Unmarshal(conf_bs, &confs)
	if err != nil {
		return nil, Utiles.NewPanicError(err)
	}

	result := []generant.GetNewer{}
	for _, conf := range confs {
		delayTime_S := conf["delayBetweenCatchRound"]

		delayTime, err := time.ParseDuration(delayTime_S)
		if err != nil {
			delayTime = _DEFAULT_DELAY
		}

		var (
			gn generant.GetNewer
		)

		switch conf["type"] {
		case "defaultChannel":
			channelName := conf["channelName"]
			if _, have := channelsDefault[channelName]; !have {
				return nil, Utiles.NewError("can't find builtin configure when type is defaultChannel")
			}
			gn = nenews_api.NewChannController(channelsDefault[channelName].Name, channelsDefault[channelName].ID, delayTime)
		case "customChannel":
			name := conf["name"]
			id := conf["id"]
			if name == "" || id == "" {
				return nil, Utiles.NewError("name or id is empty when type is customChannel")
			}
			gn = nenews_api.NewChannController(name, id, delayTime)
		case "special":
			id := conf["topicId"]
			if id == "" {
				return nil, Utiles.NewError("topicId is empty when type is special")
			}
			gn = nenews_api.NewTopicController(id, delayTime)
		default:
			return nil, Utiles.NewError("not supported type of nenews conf")
		}

		result = append(result, gn)
	}
	return result, nil
}

func init() {
	generant.Register("NeteaseNews", (generant.LoadConf)(LoadConf))
}
