package db

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

var ConsulClient *api.Client

func SetupConsul(consulAddr string) {
	ConsulClient = InitConsul(consulAddr)
}

func InitConsul(consulAddr string) *api.Client {
	consulConfig := api.DefaultConfig()
	if len(consulAddr) > 0 {
		consulConfig.Address = consulAddr
	}
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Errorf("err happened during Init consul: %s\n", err.Error())
		return nil
	}
	return consulClient
}

func ConstructUserConsulKey(owner string) string {
	return "strategy/" + owner
}

func ConstructStrategyConsulKey(owner, strategyID string) string {
	return ConstructUserConsulKey(owner) + "/" + strategyID
}

func GetConsulValue(key string) (string, error) {
	kv := ConsulClient.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		fmt.Printf("failed get %s value, error:%s\n", key, err.Error())
		return "", err
	}
	if pair == nil {
		fmt.Printf("failed get %s value, nil return\n", key)
		return "", nil
	}
	return string(pair.Value), nil
}

//Todo
func IterateConsulValues(key string) ([]string, []string, error) {
	return []string{}, []string{}, nil
}

func SetConsulValue(key string, value string) (bool, error) {
	fmt.Printf("set consule key:%s value:%s\n", key, value)
	p := &api.KVPair{Key: key, Value: []byte(value)}
	kv := ConsulClient.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		fmt.Printf("failed to set key %s value, error: %s", key, err.Error())
		return false, err
	}
	fmt.Printf("set consule key:%s value:%s\n", key, value)
	return true, nil
}

func DeleteConsulValue(key string) (bool, error) {
	kv := ConsulClient.KV()
	_, err := kv.Delete(key, nil)
	if err != nil {
		fmt.Printf("failed to delete key %s value, error: %s", key, err.Error())
		return false, err
	}
	return true, nil
}
