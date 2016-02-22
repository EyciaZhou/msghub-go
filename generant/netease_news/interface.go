package netease_news

import (
	"errors"
	"time"
	"git.eycia.me/eycia/msghub/generant"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func mustString(m map[string]interface{}, key string) (string, bool) {
	if va, hv := m[key]; hv {
		v, ok := va.(string)
		return v, ok
	}
	return "", false
}

func mustInt64Default(m map[string]interface{}, key string, defau int64) (int64) {
	if va, hv := m[key]; hv {
		v, ok := va.(float64)
		if !ok {
			return defau
		}
		return int64(v)
	}
	return defau
}

func mustDurationDefault(m map[string]interface{}, key string, defau time.Duration) time.Duration {
	if va, hv := m[key]; hv {
		v, ok := va.(string)
		if !ok {
			return defau
		}
		psd, err := time.ParseDuration(v)
		if err != nil {
			return defau
		}
		return psd
	}
	return defau
}

func loadConfigure(cf map[string]interface{}) (NeteaseNewsConfigure, error) {
	if _typ, hv := cf["type"]; hv {
		if _type, ok := _typ.(string); ok {
			if _type == "special" {
				return loadTopicConfigure(cf)
			}
		}
	}

	return loadChannelConfigure(cf)
}

func loadTopicConfigure(cf map[string]interface{}) (*NeteaseNewsTopicConfigure, error) {
	tid, can := mustString(cf, "topicId")
	if !can {
		return nil, errors.New("missing topicId field")
	}
	delayBetweenCatchRound := mustDurationDefault(cf, "delayBetweenCatchRound", time.Minute*30)

	return NewNeteaseNewsTopicConfigure(tid, delayBetweenCatchRound)

}

func loadChannelConfigure(cf map[string]interface{}) (*NeteaseNewsChannelConfigure, error) {
	if _default, hv := cf["default"]; hv {
		if _def, ok := _default.(bool); ok {
			if _def {
				if chanName, hv := mustString(cf, "channelName"); !hv {
					return nil, errors.New("use default configure but no channel name[fiele ChannelName] or type isn't string")
				} else {
					return NewDefaultNeteaseNewsChannelConfigure(chanName)
				}
			}
		} else {
			return nil, errors.New("type of fiele 'default' should be boolean")
		}
	}

	pagesOneTime := int(mustInt64Default(cf, "pagesOneTime", 2))

	delayBetweenPage := mustDurationDefault(cf, "delayBetweenPage", time.Second*10)
	delayBetweenCatchRound := mustDurationDefault(cf, "delayBetweenCatchRound", time.Minute*30)

	if chanName, hv := mustString(cf, "channelName"); !hv {
		hv := true;

		name, h := mustString(cf, "name")
		hv = hv && h

		url, h := mustString(cf, "url")
		hv = hv && h

		id, h := mustString(cf, "id")
		hv = hv && h

		if !hv {
			return nil, errors.New("channelName have wrong type, or one or some of [name, url, id] not exist in config file or the type isn't string")
		}

		return NewNeteaseNewsCatchConfigure(
			NeteaseNewsChannel{
				name, url, id,
			} ,pagesOneTime, delayBetweenPage, delayBetweenCatchRound,
		)
	} else {
		if cha, hv := channelsDefault[chanName]; !hv {
			return nil, errors.New("use builtin channel but no such channel")
		} else {
			return NewNeteaseNewsCatchConfigure(
				*cha, pagesOneTime, delayBetweenPage, delayBetweenCatchRound,
			)
		}
	}
}

func LoadGenerant(raw []byte) (generant.Generant, error) {
	var cf []map[string]interface{}

	err := json.Unmarshal(raw, &cf)

	if err != nil {
		return nil, errors.New("Netease load config error: " + err.Error())
	}

	configs := make([]NeteaseNewsConfigure, len(cf))

	for i, v := range cf {
		conf, e := loadConfigure(v)
		if e != nil {
			return nil, fmt.Errorf("Netease load config error at item %d: " + e.Error(), i+1)
		}
		configs[i] = conf
	}

	log.WithField("Plugin", "Netease News")

	return &NeteaseNewsGenerant{
		configs[:],
	}, nil
}

/*
LoadConfigAndNew:
	args should be some of *NeteaseNewsCatchConfigure, or only contain a slice of *NeteaseNewsCatchConfigure
 */
/*
func NewNeteaseNewsGenerant(args ...interface{}) (*NeteaseNewsGenerant, error) {
	if len(args) == 0 {
		return nil, errors.New("at least one config to catch")
	}

	typeOfConfig := reflect.TypeOf(&NeteaseNewsCatchConfigure{})

	if len(args) == 1 {
		if reflect.TypeOf(args[0]) == reflect.SliceOf(typeOfConfig) {
			return &NeteaseNewsGenerant{
				args[0].([]*NeteaseNewsCatchConfigure),
			}, nil
		}
	}

	configs := make([]*NeteaseNewsCatchConfigure, len(args))

	for i := 0; i < len(args); i++ {
		if reflect.TypeOf(args[i]) != typeOfConfig {
			return nil, errors.New("config type should be *NeteaseNewsCatchConfigure")
		}
		configs[i] = args[i].(*NeteaseNewsCatchConfigure)
	}

	return &NeteaseNewsGenerant{
		configs[:],
	}, nil
}
*/

type NeteaseNewsGenerant struct {
	configs []NeteaseNewsConfigure
}

func (n *NeteaseNewsGenerant) Catch() {
	for _, config := range n.configs {
		go catchDaemon(config)
	}
}

func (n *NeteaseNewsGenerant) Stop() {
	for _, conf := range n.configs {
		stopCatchDaemon(conf)
	}
}

func (n *NeteaseNewsGenerant) ForceStop() {
	n.Stop()
}

func init() {
	generant.Register("NeteaseNews", generant.LoadConf(LoadGenerant))
}