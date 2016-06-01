package plugins

import (
	"github.com/EyciaZhou/configparser"
	log "github.com/Sirupsen/logrus"
	"time"
)

type config_t struct {
	ConfDir         string
	ConfFileNames   []string
	ConfPluginNames []string
}

var (
	config config_t
)

func loadConfig() {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	configparser.AutoLoadConfig("plugins", &config)
	configparser.ToJson(&config)
}

func Init() {
	loadConfig()
	log.Info("Start load plugins's config")
	err := loadPluginConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	log.Infof("Start fire plugins, %d plugins to fire", len(pluginRunners))
	for i, plugin := range pluginRunners {
		log.Infof("[%d/%d]...", i+1, len(pluginRunners))
		go plugin.Catch()
		log.Info("fired and start delay")
		time.Sleep(10 * time.Second)
	}

	log.Info("Init finished")
}
