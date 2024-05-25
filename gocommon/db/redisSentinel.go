package db

import (
	"fmt"

	"github.com/go-redis/redis"
)

func SetupRedisSentinel(masterName string, sentinelAddresses []string, password string, db int) {
	RClient = InitRedisSentinel(masterName, sentinelAddresses, password, db)
}

func InitRedisSentinel(masterName string, sentinelAddresses []string, password string, db int) *redis.Client {
	RClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAddresses,
		Password:      password,
		DB:            db,
	})
	_, err := RClient.Ping().Result()
	if err != nil {
		fmt.Printf("redis connection err: %v\n", err)
	} else {
		fmt.Printf("redis connection succeed\n")
	}
	return RClient
}
