package commoncontext

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"request-matcher-openai/gocommon/db"
)

func setupRedis(vp *viper.Viper) (*redis.Client, error) {
	if vp.IsSet("redis.host") {
		rclient := db.LoadAndSetupRedis(vp)
		if rclient == nil {
			return nil, fmt.Errorf("failed to setup redis client")
		}
		return rclient, nil
	} else {
		GetInstance().MyLogger.Infof("redis host is not set, so skip to setup redis client")
		return nil, nil
	}
}
