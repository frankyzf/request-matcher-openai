package commoncontext

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"request-matcher-openai/gocommon/viperconf"
)

func GetAndMergeConfig(vp *viper.Viper, env string) *viper.Viper {
	fmt.Printf("merge config with env:%v\n", env)

	viperconf.LoadConfiguration(env, vp)
	vp.MergeInConfig()
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		vp = viperconf.LoadConfiguration(env, vp)
		vp.MergeInConfig()
		fmt.Printf("Config file changed::%v\n", e.Name)
		viperconf.DumpConfiguration(vp, "./config.dump.json")
	})
	viperconf.DumpConfiguration(vp, "./config.dump.json")
	return vp
}
