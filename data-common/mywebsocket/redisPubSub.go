package mywebsocket

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/commoncontext"
)

type RedisPubSub struct {
	userSession *UserSession
	Connection  *websocket.Conn
	myRClient   *redis.Client
	Running     bool
	mylogger    *logrus.Entry
}

func GetRedisPubSub(session *UserSession, conn *websocket.Conn, rclient *redis.Client) (*RedisPubSub, error) {
	p := &RedisPubSub{
		userSession: session,
		Connection:  conn,
		myRClient:   rclient,
		Running:     false,
		mylogger:    commoncontext.SetupLogging("redis", "pubsub"),
	}
	p.mylogger.Infof("start redis pubsub")
	p.clearQueue()
	return p, nil
}

func (p *RedisPubSub) Start() {
	p.mylogger.Infof("RedisPubSub start!!")
	p.Running = true
	go func() {
		for p.Running {
			redisKeys := p.allRedisKeys()
			bProcess := false
			p.mylogger.Debugf("redispub redisKeys:%v", redisKeys)
			for _, redisKey := range redisKeys {
				num, _ := p.myRClient.LLen(redisKey).Result()
				p.mylogger.Debugf("redispub redisKey:%v, num:%v", redisKeys, num)
				if num > 0 {
					bProcess = true
					notificationJsonData, err2 := p.myRClient.LPop(redisKey).Result()
					if err2 != nil {
						p.mylogger.Errorf("failed to LPop %v, err:%v", redisKey, err2)
					} else {
						notification := export.WSNotification{}
						json.Unmarshal([]byte(notificationJsonData), &notification)
						topic := getTopicFromKey(redisKey)
						p.mylogger.Infof("redispub get a topic:%v,notification:%v", topic, notification)
						subItem, err := GetSubItemByTopic(topic)
						if err != nil {
							p.mylogger.Errorf("failed to get subitem for topic:%v", topic)
							continue
						}
						// todo
						go func() {
							buf, err3 := subItem.GetSnapshot(topic, notification)
							if err3 != nil || len(buf) == 0 {
								p.mylogger.Warningf("skip to send packet because failed to get snapshot:%v", err3)
							} else {
								ConnectionSendData(p.Connection, []byte(buf))
								p.mylogger.Debugf("redispub sending client:%v with a packet:%v", redisKey, string(buf))
							}
						}()
					}
				}
			}
			if bProcess == false {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (p *RedisPubSub) Stop() {
	p.mylogger.Infof("RedisPubSub Stop!!")
	p.clearQueue()
	p.Running = false
}

func (p *RedisPubSub) Publish(redisKey string, buf []byte) error {
	p.mylogger.Infof("redis publish data for key:%v ", redisKey)
	topic := getTopicFromKey(redisKey)

	num, _ := p.myRClient.LLen(redisKey).Result()
	if num > int64(commoncontext.GetDefaultInt("websocket.redis_queue_length", 1)) {
		p.mylogger.Infof("skip to publish redis:%v because the notification is there already", redisKey)
	} else {
		_, err := p.myRClient.RPush(redisKey, buf).Result()
		if err == nil {
			p.mylogger.Infof("successfully publish topic:%v to redisKeys:%v", topic, redisKey)
		} else {
			p.mylogger.Errorf("error to publish for key:%v, err:%v", redisKey, err)
		}
	}

	return nil
}

func (p *RedisPubSub) OnSubscribe(topic string, session *UserSession) error {
	p.mylogger.Infof("redisPubsub subscribe topic :%v", topic)
	subItem, err := GetSubItemByTopic(topic)
	if err != nil {
		p.mylogger.Errorf("no matched subItem for topic:%v", topic)
		return err
	}
	if subItem.SendDataWhenSubscribe(topic) {
		go func() {
			data, err := subItem.GetSnapshot(topic, export.WSNotification{})
			if len(data) == 0 || err != nil {
				p.mylogger.Errorf("skip to send empty snapshot for the first time, err:%v", err)
			} else {
				err = ConnectionSendData(p.Connection, data)
				// _, err := p.myRClient.RPush(redisKey, data).Result()
				if err != nil {
					p.mylogger.Errorf("error to send data err:%v", err)
				} else {
					p.mylogger.Infof("successfully send snapshot for topic:%v and id:%v", topic, session.ID)
				}
			}
		}()
	} else {
		p.mylogger.Warningf("skip to send snapshot when subscribe:%v", topic)
	}
	return nil
}

func (p *RedisPubSub) OnUnsubscribe(topic string, session *UserSession) error {
	p.mylogger.Infof("redisPubsub unsubscribe topic :%v", topic)
	return nil
}

func (p *RedisPubSub) OnChangeSubscribe(oldTopic string, newTopic string, session *UserSession) error {
	p.mylogger.Infof("redisPubsub change subscribe from topic :%v to new topic:%v", oldTopic, newTopic)
	return nil
}

func (p *RedisPubSub) clearQueue() error {
	regKeys := getKey("*", p.userSession.ID)
	keys, _ := p.myRClient.Keys(regKeys).Result()
	p.mylogger.Infof("clearQueue:%v with regkeys:%v", keys, regKeys)
	for _, key := range keys {
		p.myRClient.Del(key)
		p.mylogger.Debugf("delete redis key:%v", key)
	}
	return nil
}

func (p *RedisPubSub) allRedisKeys() []string {
	regKeys := getKey("*", p.userSession.ID)
	keys, _ := p.myRClient.Keys(regKeys).Result()
	return keys
}
