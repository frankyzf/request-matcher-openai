package mywebsocket

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"request-matcher-openai/data-common/datamanage"
	"request-matcher-openai/data-common/myauth"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/commoncontext"
)

var mylogger *logrus.Entry
var manage *datamanage.DataManager
var auth *myauth.MyAuth
var rclient *redis.Client
var connMu sync.Mutex

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Setup(dm *datamanage.DataManager, mauth *myauth.MyAuth) (export.WebSocketManagerInterface, error) {
	manage = dm
	auth = mauth
	connMu = sync.Mutex{}
	var p export.WebSocketManagerInterface
	mylogger = commoncontext.SetupLogging("0", "websocket")
	mytype := commoncontext.GetDefaultString("websocket.type", "redis")
	rclient = commoncontext.GetInstance().RClient
	if rclient == nil {
		return nil, errors.New("failed to setup redis")
	}
	p = GetWebsocketManager(mytype, rclient, mauth)
	return p, nil
}

func ConnectionSendData(conn *websocket.Conn, body []byte) error {
	connMu.Lock()
	defer connMu.Unlock()
	if err := conn.WriteMessage(websocket.TextMessage, body); err == nil {
		mylogger.Debugf("websocket sendout msg:%v", string(body))
		return nil
	} else {
		mylogger.Errorf("websocket failed to send packet, err:%v", err)
		return err
	}
}

func getKey(topic string, clientSubID string) string {
	return fmt.Sprintf("pubsub:%v:%v", topic, clientSubID)
}
