package netease_news

import (
	"reflect"
	"errors"
	"time"
	"git.eycia.me/eycia/msghub/generant"
	"encoding/json"
)

type NeteaseNewsCatchConfigure struct {
	NeteaseNewsChannel
	PagesOneTime		int
	DelayBetweenPage	time.Duration
	DelayBetweenEachCatchRound time.Duration

	quit chan int
}

func LoadConf(raw []byte) (generant.Generant, error) {
	var cf map[string]interface{}

	json.Unmarshal(raw, &cf)
}

func NewDefaultNeteaseNewsCatchConfigure(channelName string) (*NeteaseNewsCatchConfigure, error) {
	if _, hv := channelsDefault[channelName]; !hv {
		return nil, errors.New("no such channel")
	}
	channelInfo := channelsDefault[channelName]
	return &NeteaseNewsCatchConfigure{
		*channelInfo,
		2,
		time.Second * 10,
		time.Minute * 10,

		make(chan int),
	}, nil
}

func NewNeteaseNewsCatchConfigure(chann NeteaseNewsChannel, pagesOneTime int,
delayBetweenPage, delayBetweenEachCatchRound time.Duration) (*NeteaseNewsCatchConfigure, error) {
	return &NeteaseNewsCatchConfigure{
		chann,
		pagesOneTime,
		delayBetweenPage,
		delayBetweenEachCatchRound,

		make(chan int),
	}, nil
}

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

func init() {
	generant.Register("NeteaseNews", generant.LoadConf(LoadConf))
}