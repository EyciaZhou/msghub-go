package netease_news

import (
	"reflect"
	"errors"
)

/*
LoadConfigAndNew:
	args should be some of *NeteaseNewsCatchConfigure, or only contain a slice of *NeteaseNewsCatchConfigure
 */
func NewNeteaseNewsGenerant(args ...interface{}) (*NeteaseNewsGenerant, error) {
	if len(args) == 0 {
		return nil, errors.New("at least one config to catch")
	}

	typeOfConfig := reflect.TypeOf(&NeteaseNewsCatchConfigure{})

	if len(args) == 1 {
		if reflect.TypeOf(args[0]) == reflect.SliceOf(typeOfConfig) {
			return NeteaseNewsCatchConfigure{
				args[0].([]*NeteaseNewsCatchConfigure),
			}, nil
		}
	}

	configs := make(NeteaseNewsCatchConfigure, len(args))

	for i := 0; i < len(args); i++ {
		if reflect.TypeOf(args[i]) != typeOfConfig {
			return nil, errors.New("config type should be *NeteaseNewsCatchConfigure")
		}
		configs[i] = args[i].(*NeteaseNewsCatchConfigure)
	}

	return NeteaseNewsCatchConfigure{
		configs[:],
	}, nil
}

type NeteaseNewsGenerant struct {
	configs []*NeteaseNewsCatchConfigure
}

func (n *NeteaseNewsGenerant) Catch() {
	for _, config := range n.configs {
		go config.CatchDaemon()
	}
}

func (n *NeteaseNewsGenerant) Stop() {
	for _, conf := range n.configs {
		conf.Stop()
	}
}

func (n *NeteaseNewsGenerant) ForceStop() {
	n.Stop()
}