package mywebsocket

import (
	"errors"
	"request-matcher-openai/data-mydb/export"
)

type SubItemInterface interface {
	GetName() string
	GetSubscribeDataType() []string
	GetSnapshot(topic string, notification export.WSNotification) ([]byte, error)
	SendDataWhenSubscribe(topic string) bool
}

func GetSubItemByTopic(topic string) (SubItemInterface, error) {
	if topic == "xxx" {
		return nil, nil
	} else {
		return nil, errors.New("failed to get the subitem for topic:" + topic)
	}
}

func GetSubscribeDataTypesByTopic(topic string) []string {
	if topic == "xxx" {
		return []string{}
	} else {
		mylogger.Warningf("unknown subscribe types")
		return []string{}
	}
}
