package plugins

import (
	//"errors"
	//"fmt"
	"github.com/EyciaZhou/msghub.go/interface"
	//log "github.com/Sirupsen/logrus"
	//"io/ioutil"
	//"path"
	"sync"
	"time"
)

type Plugin interface {
	ResumeStatus(Status []byte) (PluginTask, error)
	LoadConfig(Config interface{}) (PluginTask, error)
	GetConfigType() interface{}
}

type PluginTask interface {
	FetchNew() (*Interface.Topic, error)

	GetDelayBetweenCatchRound() time.Duration

	DumpStatus() (Status []byte)

	Hash() []byte //as id, less equal 64 bytes
}

var (
	pluginsMu       sync.Mutex
	RegistedPlugins = make(map[string]Plugin)

	pluginRunners []*pluginRunner
)

func Register(name string, plugin Plugin) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	RegistedPlugins[name] = plugin
}

func loadPluginConfig() error {
	//
	/*
		if pluginRunners != nil && len(pluginRunners) > 0 {
			return errors.New("can't load config twice")
		}

		if len(config.ConfFileNames) != len(config.ConfPluginNames) {
			return errors.New("the length of ConfFileNames and ConfPluginNames not same")
		}

		pluginRunners = []*pluginRunner{}

		config.ConfDir = path.Clean(config.ConfDir)

		log.Infof("%d plugin configs to load", len(config.ConfFileNames))

		for i, fn := range config.ConfFileNames {
			log.Infof("[%d/%d]...", i+1, len(config.ConfFileNames))
			bs, err := ioutil.ReadFile(config.ConfDir + "/" + fn)
			if err != nil {
				return err
			}

			if plugin, hv := RegistedPlugins[config.ConfPluginNames[i]]; !hv {
				return errors.New("Doesn't registed Plugin : " + config.ConfPluginNames[i])
			} else {
				gns, err := plugin.LoadConf(bs)
				if err != nil {
					return fmt.Errorf("Error when load %d th plugin : %s", i+1, err.Error())
				}
				for _, gn := range gns {
					pluginRunners = append(pluginRunners, NewPluginRunner(gn))
				}
			}
		}

		return nil
	*/
	return nil
}
