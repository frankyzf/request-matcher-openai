package export

import (
	"github.com/gin-gonic/gin"
)

type WSNotification struct {
	DataType string      `json:"data_type"`
	Data     interface{} `json:"data"`
}

type WebSocketManagerInterface interface {
	AcceptSession(c *gin.Context) error
	SendData(dataType string, notification WSNotification)
	Start() error
	Stop() error
	GetName() string
	OnDataChangedManually(dataType string, notification WSNotification)
}
