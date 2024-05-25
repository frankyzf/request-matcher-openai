package db

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"request-matcher-openai/gocommon/viperconf"
)

var RedisHost, RedisPort, RedisPassword string
var RedisDB int
var RedisURL string

func LoadAndSetupRedis(vp *viper.Viper) *redis.Client {
	connectionMode := viperconf.GetDefaultString(vp, "redis.connection_mode", "single")
	if connectionMode == "sentinel" {
		return loadSentinelRedis(vp)
	}
	return loadSingleRedisInstance(vp)
}

func loadSingleRedisInstance(vp *viper.Viper) *redis.Client {
	RedisHost = viperconf.GetDefaultString(vp, "redis.host", "localhost")
	RedisPort = viperconf.GetDefaultString(vp, "redis.port", "6379")
	RedisPassword = viperconf.GetDefaultString(vp, "redis.password", "")
	RedisDB = viperconf.GetDefaultInt(vp, "redis.db", 0)
	RedisURL = fmt.Sprintf("%s:%s", RedisHost, RedisPort)
	conn := InitRedis(RedisURL, RedisPassword, RedisDB)
	return conn
}

func loadSentinelRedis(vp *viper.Viper) *redis.Client {
	masterName := viperconf.GetDefaultString(vp, "redis.master_name", "")
	sentinelAddresses := viperconf.GetDefaultStringSlice(vp, "redis.sentinel_addresses", []string{})
	RedisPassword = viperconf.GetDefaultString(vp, "redis.password", "")
	RedisDB = viperconf.GetDefaultInt(vp, "redis.db", 0)
	if masterName == "" {
		fmt.Printf("sentinel empty master name\n")
		return nil
	}
	if len(sentinelAddresses) == 0 {
		fmt.Printf("sentinel empty sentinel addresses\n")
		return nil
	}
	conn := InitRedisSentinel(masterName, sentinelAddresses, RedisPassword, RedisDB)
	return conn
}
