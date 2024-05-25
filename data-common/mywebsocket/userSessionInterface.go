package mywebsocket

import (
	"github.com/gorilla/websocket"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/data-mydb/mydb"
)

type UserSessionInterface interface {
	GetID() string
	GetOwner() mydb.BaseAccount
	SetOwner(mydb.BaseAccount)
	GetConnection() *websocket.Conn
	Terminate()
	Start() error
	OnUpdatedData(dataType string, notification export.WSNotification) error
	IsExpired() bool
}

type DataChangeTrackerInterface interface {
	GetID() string
	OnUpdatedData(dataType string, notification export.WSNotification) error
}
