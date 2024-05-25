package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var RClient *redis.Client

func SetupRedis(url string, password string, db int) {
	RClient = InitRedis(url, password, db)
}

func InitRedis(url string, password string, db int) *redis.Client {
	RClient := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	_, err := RClient.Ping().Result()
	if err != nil {
		fmt.Printf("redis connction err: %v\n", err)
	}
	return RClient
}

func ConstructStatusKey(runnerId string) string {
	return "Status:" + runnerId
}

func ConstructUserStatusKey(user string) string {
	return "Status:" + user
}

func GetStatusValue(key string) (string, error) {
	val, err := RClient.Get(key).Result()
	fmt.Printf("get redis-key: %s, val:%s\n", key, val)
	if err == nil {
		return val, err
	} else if err == redis.Nil {
		fmt.Printf("redis-key:%s does not exist\n", key)
		return "", nil
	} else {
		return "", err
	}
}

func GetUserAllStatus(uid string) (map[string]string, error) {
	k := ConstructUserStatusKey(uid) + ":*"
	return IterateStatusValues(k)
}

func IterateStatusValues(regPath string) (map[string]string, error) {
	var cursor uint64
	var resKeys = []string{}
	var res = make(map[string]string)
	var err error
	for {
		var keys []string
		keys, cursor, err = RClient.Scan(cursor, regPath, 20).Result()
		if err != nil {
			break
		}
		resKeys = append(resKeys, keys...)
		if cursor == 0 {
			break
		}
	}
	var v string
	for _, k := range resKeys {
		v, err = GetStatusValue(k)
		if err != nil {
			break
		}
		res[k] = v
	}
	fmt.Printf("iterate status : %v \n", res)
	return res, err
}

func SetStatusValue(key string, status string) (bool, error) {
	fmt.Printf("set %s status to %s\n", key, status)
	err := RClient.Set(key, status, 0).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func SetStatusValueTimeout(key string, status string, second int64) (bool, error) {
	fmt.Printf("set %s status to %s with timeout:%d\n", key, status, second)
	err := RClient.Set(key, status, time.Duration(second)*time.Second).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
