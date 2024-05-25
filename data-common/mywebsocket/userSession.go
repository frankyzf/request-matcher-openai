package mywebsocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/util"
)

type UserSession struct {
	ID             string
	Owner          mydb.BaseAccount
	wsManager      *WebSocketManager
	Connection     *websocket.Conn
	PubSub         *RedisPubSub
	pubsubInit     bool
	SubscriberKeys sync.Map // key(subtopic_subscriberID) => data_type
	RClient        *redis.Client
	LastUpdate     time.Time
	Running        bool
	mylogger       *logrus.Entry
}

func MakeUserSession(wsmanager *WebSocketManager,
	conn *websocket.Conn, rclient *redis.Client) UserSessionInterface {
	p := &UserSession{
		ID:             util.GetUUID(),
		wsManager:      wsmanager,
		Connection:     conn,
		RClient:        rclient,
		pubsubInit:     false,
		SubscriberKeys: sync.Map{},
		Running:        false,
		mylogger:       commoncontext.SetupLogging("session", "websocket"),
	}
	return p
}

func (p *UserSession) GetID() string {
	return p.ID
}

func (p *UserSession) SetOwner(owner mydb.BaseAccount) {
	p.Owner = owner
}

func (p *UserSession) GetOwner() mydb.BaseAccount {
	return p.Owner
}

func (p *UserSession) GetConnection() *websocket.Conn {
	return p.Connection
}

func (p *UserSession) IsExpired() bool {
	expire := commoncontext.GetDefaultInt("websocket.client_expire", 300)
	now := time.Now()
	if now.Sub(p.LastUpdate).Seconds() > float64(expire) {
		p.mylogger.Warningf("the session:%v expired", p.ID)
		return true
	}
	return false
}

func (p *UserSession) Start() error {
	p.mylogger.Infof("user session start")
	var err error
	p.updateTimestamp()
	p.PubSub, err = GetRedisPubSub(p, p.Connection, p.RClient)
	if err != nil {
		p.mylogger.Errorf("failed to start user session:%v", err)
		return err
	}
	p.PubSub.Start()
	p.pubsubInit = true
	p.Running = true
	go func() {
		for p.Running {
			messageType, buf, err := p.Connection.ReadMessage()
			if err != nil {
				p.mylogger.Errorf("error to read the buf data, err:%v", err)
				p.wsManager.RemoveSession(p.ID) //the termination will called in the RemoveSession()
				return
			}
			if messageType == websocket.TextMessage {
				err = p.handleMessage(buf)
				if err != nil {
					p.mylogger.Errorf("failed to handle message:%v and err:%v", string(buf), err)
				} else {
					p.mylogger.Debugf("successfully handle the messsage:%v", string(buf))
				}
			} else {
				p.mylogger.Warningf("skip to handle a non text message type:%v, content:%v", messageType, string(buf))
			}
		}
	}()
	return nil
}

func (p *UserSession) Terminate() {
	p.mylogger.Warningf("user serssion:%v is terminated", p.ID)
	p.Running = false
	if p.pubsubInit {
		p.PubSub.Stop()
	}
	p.Connection.Close()
	return
}

// redisPubSub insert a notification flag to the redis queue
func (p *UserSession) OnUpdatedData(dataType string, notification export.WSNotification) error {
	keys := p.getSubscribeKeysByDataType(dataType)
	p.mylogger.Infof("session ready to publish for dataType:%v and keys:%v ", dataType, keys)

	buf, _ := json.Marshal(notification)
	for _, key := range keys {
		p.PubSub.Publish(key, buf)
	}
	return nil
}

func (p *UserSession) handleMessage(rawMsg []byte) error {
	msg := wsMessage{}
	err := json.Unmarshal(rawMsg, &msg)
	p.mylogger.Debugf("recv websocket msg:%v", string(rawMsg))
	if err != nil {
		p.mylogger.Errorf("error in unmarshal message:%v, err:%v", rawMsg, err)
		return err
	}
	if msg.OpType == "subscribe" {
		p.subscribe(msg)
		return nil
	} else if msg.OpType == "unsubscribe" {
		p.unsubscribe(msg)
		return nil
	} else if msg.OpType == "keepAlive" {
		buf, _ := json.Marshal(msg)
		err = ConnectionSendData(p.Connection, buf)
		p.updateTimestamp()
	} else {
		p.mylogger.Errorf("recv unknown message type:%v, topic:%v", msg.OpType, msg.Data)
		return errors.New("unknown message type:" + msg.OpType)
	}
	return nil
}

func (p *UserSession) subscribe(msg wsMessage) {
	p.mylogger.Infof("subscribe topic:%v", msg.Data)
	topic := p.getTopicForSubscribe(msg)
	if topic == "" {
		subReply := wsMessage{OpType: "subscribe_reply", Data: "invalid topic"}
		buf, _ := json.Marshal(subReply)
		ConnectionSendData(p.Connection, buf) //reply first
	} else {
		subReply := wsMessage{OpType: "subscribe_reply", Data: "success"}
		buf, _ := json.Marshal(subReply)
		ConnectionSendData(p.Connection, buf) //reply first

		key := getKey(topic, p.ID)
		p.addDataTopic(topic, key)
		//send a snapshot immediately after subscribe
		p.PubSub.OnSubscribe(topic, p)
	}
}

func (p *UserSession) getTopicForSubscribe(msg wsMessage) string {
	if msg.Data == nil {
		return ""
	}
	if data, ok := msg.Data.(string); ok {
		return data
	}
	return ""
}

func (p *UserSession) unsubscribe(msg wsMessage) {
	p.mylogger.Infof("unsubscribe topic:%v session:%v", msg.Data, p.ID)
	topic := p.getTopicForSubscribe(msg)
	buf, _ := json.Marshal(wsMessage{OpType: "unsubscribe_reply", Data: "success"})
	ConnectionSendData(p.Connection, buf)

	key := getKey(topic, p.ID)
	p.removeDataTopic(topic, key)
	p.PubSub.OnUnsubscribe(topic, p)
}

func (p *UserSession) addDataTopic(topic string, key string) {
	p.mylogger.Infof("add data topic:%v key:%v for session:%v", topic, key, p.ID)
	dataTypes := GetSubscribeDataTypesByTopic(topic)
	for _, dataType := range dataTypes {
		p.mylogger.Infof("add data type:%v for topic:%v and  session:%v", dataType, topic, p.ID)
		p.addSubscribeKey(key, dataType)
		p.wsManager.AddDataTypeSession(dataType, p)
	}
}

func (p *UserSession) removeDataTopic(topic string, key string) {
	p.mylogger.Infof("add data topic:%v key:%v for session:%v", topic, key, p.ID)
	dataTypes := GetSubscribeDataTypesByTopic(topic)
	for _, dataType := range dataTypes {
		p.mylogger.Infof("remove data type:%v for topic:%v and  session:%v", dataType, topic, p.ID)
		p.removeSubscribeKey(dataType)
		count := 0
		p.SubscriberKeys.Range(func(keyT, value interface{}) bool {
			mDataType, _ := keyT.(string)
			if mDataType == dataType {
				count++
			}
			return true
		})
		if count <= 1 {
			p.wsManager.RemoveDataTypeSession(dataType, p)
		}
	}
}

func (p *UserSession) addSubscribeKey(key string, dataType string) {
	p.mylogger.Infof("successfully add subscriber with key:%v, dataType:%v with session:%v", key, dataType, p.ID)
	p.SubscriberKeys.Store(dataType, key)
}

func (p *UserSession) removeSubscribeKey(key string) {
	p.mylogger.Infof("delete subscribe key:%v ", key)
	if _, ok := p.SubscriberKeys.Load(key); ok {
		p.SubscriberKeys.Delete(key)
	}
}

func (p *UserSession) getSubscribeKeysByDataType(dataTypeV string) []string {
	res := []string{}
	p.SubscriberKeys.Range(func(keyBuf, value interface{}) bool {
		dataType, _ := keyBuf.(string)
		if dataType == dataTypeV {
			key, _ := value.(string)
			res = append(res, key)
		}
		return true
	})
	return res
}

func (p *UserSession) updateTimestamp() {
	p.LastUpdate = time.Now()
	p.mylogger.Infof("user session:%v late update is set to:%v", p.ID, p.LastUpdate)
}

func (p *UserSession) findKeyBySessionID(sessionID string) []string {
	data := []string{}
	subfix := fmt.Sprintf(":%v", sessionID)
	p.SubscriberKeys.Range(func(keyBuf, value interface{}) bool {
		key := value.(string)
		if strings.HasSuffix(key, subfix) {
			data = append(data, key)
		}
		return true
	})
	return data
}

func getTopicFromKey(key string) string {
	dd := strings.Split(key, ":")
	if len(dd) != 3 {
		mylogger.Errorf("the key is not correct format:%v", key)
		return ""
	}
	return dd[1]
}
