package config

import (
	"fmt"
	"path"

	"github.com/Ferriem/Distributed_schedule_platform/common/models"
	"github.com/Ferriem/Distributed_schedule_platform/common/pkg/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	ExtensionJson = ".json"
	ExtensionYaml = ".yaml"
	ExtensionIni  = ".ini"
	NameSpace     = "conf"
)

var (
	autoLoadLocalConfigs = []string{
		ExtensionJson,
		ExtensionYaml,
		ExtensionIni,
	}
)

var _defaultConfig *models.Config

func LoadConfig(env, serverName, configFileName string) (*models.Config, error) {
	var c models.Config
	var confPath string
	dir := fmt.Sprintf("%s/%s/%s", serverName, NameSpace, env)
	for _, registerExt := range autoLoadLocalConfigs {
		confPath = path.Join(dir, configFileName+registerExt)
		if utils.Exists(confPath) {
			break
		}
	}
	fmt.Println("the path to the configuration file you are using is :", confPath)
	v := viper.New()
	v.SetConfigFile(confPath)
	ext := utils.Ext(confPath)
	v.SetConfigType(ext)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&c); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&c); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("load config is :#%v\n", c)
	_defaultConfig = &c
	return &c, nil
}

func GetConfigModels() *models.Config {
	return _defaultConfig
}
