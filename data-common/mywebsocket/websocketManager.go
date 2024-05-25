package mywebsocket

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"request-matcher-openai/data-common/myauth"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/commoncontext"
)

type WebSocketManager struct {
	mType       string
	myRClient   *redis.Client
	myAuth      *myauth.MyAuth
	Sessions    sync.Map // session_id ==> *userSession
	DataTopics  sync.Map // data_types => []DataChangeTrackerInterface (acctually it is []*userSession)
	loggerLevel string
	mylogger    *logrus.Entry
}

func GetWebsocketManager(mtype string, rclient *redis.Client, mauth *myauth.MyAuth) *WebSocketManager {
	p := &WebSocketManager{
		mType:      mtype,
		myRClient:  rclient,
		myAuth:     mauth,
		Sessions:   sync.Map{},
		DataTopics: sync.Map{},
		mylogger:   commoncontext.SetupLogging("manager", "websocket"),
	}

	p.mylogger.Infof("start a %v web socket manager", p.mType)
	return p
}

func (p *WebSocketManager) GetName() string {
	return "websocket manager"
}

func (p *WebSocketManager) AcceptSession(c *gin.Context) error {
	w := c.Writer
	r := c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		var session UserSessionInterface
		if p.mType == "redis" {
			session = MakeUserSession(p, conn, p.myRClient)
		} else {
			p.mylogger.Errorf("unknown websocket type:" + p.mType)
			return errors.New("wrong websocket type:" + p.mType)
		}
		p.mylogger.Infof("accept a user session:%v", session.GetID())
		timeout := commoncontext.GetDefaultInt("websocket.strange_expire", 300) //5 minute
		go func() {
			for {
				select {
				case err := <-p.waitingForSignature(session):
					if err != nil {
						p.mylogger.Errorf("failed to authorized:%v, err:%v", session.GetID(), err)
						buf, _ := json.Marshal(wsMessage{OpType: "signature_reply", Data: "auth failed and terminate"})
						ConnectionSendData(conn, buf)
						session.Terminate() //terminate will close connection
						return
					} else {
						p.mylogger.Infof("successfully authorized id:%v, name:%v", session.GetID(), session.GetOwner().Name)
						buf, _ := json.Marshal(wsMessage{OpType: "signature_reply", Data: "success"})
						err = ConnectionSendData(conn, buf)
						if err == nil {
							session.Start()
							p.Sessions.Store(session.GetID(), session)
						} else {
							p.mylogger.Warningf("close connection:%v due to connection send err:%v", session.GetID(), err)
							session.Terminate()
						}
						return
					}
				case <-time.After(time.Duration(timeout) * time.Second):
					p.mylogger.Errorf("expire for websocket waiting signature, so quit it:%v", session.GetID())
					session.Terminate()
					return
				}
			}
		}()
		return nil
	} else {
		p.mylogger.Errorf("error to upgrade system:%v", err)
	}
	return err
}

func (p *WebSocketManager) waitingForSignature(session UserSessionInterface) <-chan error {
	out := make(chan error)
	p.mylogger.Infof("start to wait for signature from websocket client:%v", session.GetID())
	go func() {
		messageType, buf, err := session.GetConnection().ReadMessage()
		if err != nil {
			p.mylogger.Errorf("error to read the buf data, err:%v", err)
			out <- err
			return
		}
		if messageType != websocket.TextMessage {
			p.mylogger.Warningf("recv a non text messagemessage type:%v, content:%v", messageType, string(buf))
			out <- errors.New("non text message")
			return
		}
		msg := wsMessage{}
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			out <- errors.New("failed to unmarshal signature message")
			return
		}
		if msg.OpType != "signature" {
			out <- errors.New("non signature op_type:" + msg.OpType)
			return
		}

		if msg.Data == nil {
			p.mylogger.Errorf("error to signature websocket signatureValue is nil")
			out <- errors.New("signature error")
			return
		}
		if _, ok := msg.Data.(map[string]interface{}); ok {

		} else {
			p.mylogger.Errorf("error to signature websocket signatureValue is nil")
			out <- errors.New("signature error")
		}

		out <- nil
	}()
	return out
}

func (p *WebSocketManager) SendData(dataType string, notification export.WSNotification) {
	p.sendTableUpdateToMatchedSession(dataType, notification)
}

func (p *WebSocketManager) sendTableUpdateToMatchedSession(dataType string, notification export.WSNotification) {
	sessions := p.getDataTypeSessions(dataType)
	p.mylogger.Infof("send table update to session with sessions:%v, dataType:%v", sessions, dataType)
	for _, session := range sessions {
		session.OnUpdatedData(dataType, notification)
	}
}

func (p *WebSocketManager) Start() error {
	p.mylogger.Infof("websocket manager start")
	p.monitorProcess()
	return nil
}

func (p *WebSocketManager) Stop() error {
	p.mylogger.Infof("websocket manager stop")
	return nil
}

func (p *WebSocketManager) OnDataChangedManually(dataType string, notification export.WSNotification) {
	p.mylogger.Infof("OnDataChangedManually:%v, notification:%v", dataType, notification)
	sessions := p.getDataTypeSessions(dataType)
	if len(sessions) > 0 {
		go func() {
			p.sendTableUpdateToMatchedSession(dataType, notification)
		}()
	} else {
		//do nothing
		p.mylogger.Warningf("the data of table:%v is not needed", dataType)
	}
}

func (p *WebSocketManager) monitorProcess() {
	p.mylogger.Infof("RedisPubSub monitor Process start to check nonactive session")
	go func() {
		for true {
			nonActive := map[string]UserSessionInterface{}
			p.Sessions.Range(func(key, value interface{}) bool {
				session := value.(UserSessionInterface)
				if session.IsExpired() {
					nonActive[session.GetID()] = session
				}
				return true
			})
			for _, session := range nonActive {
				p.mylogger.Warningf("kickout the non-active sessionID:%v", session.GetID())
				p.RemoveSession(session.GetID())
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (p *WebSocketManager) getDataTypeSessions(dataType string) []DataChangeTrackerInterface {
	res := []DataChangeTrackerInterface{}
	if keyBuf, ok := p.DataTopics.Load(dataType); ok {
		res, _ = keyBuf.([]DataChangeTrackerInterface)
	}
	return res
}

func (p *WebSocketManager) AddDataTypeSession(dataType string, us DataChangeTrackerInterface) {
	if buf, ok2 := p.DataTopics.Load(dataType); ok2 {
		userSessions := buf.([]DataChangeTrackerInterface)
		for _, session := range userSessions {
			if session.GetID() == us.GetID() {
				p.mylogger.Warningf("skip to add session:%v which is already add for data type:%v", us.GetID(), dataType)
				return
			}
		}
		userSessions = append(userSessions, us)
		p.DataTopics.Store(dataType, userSessions)
		p.mylogger.Infof("data topics store a dataType:%v with len(UserSession):%v", dataType, len(userSessions))
	} else {
		sessions := []DataChangeTrackerInterface{us}
		p.DataTopics.Store(dataType, sessions)
		p.mylogger.Infof("data topics store a new dataType:%v for session:%v", dataType, us.GetID())
	}
}

func (p *WebSocketManager) RemoveDataTypeSession(dataType string, us DataChangeTrackerInterface) {
	if buf, ok2 := p.DataTopics.Load(dataType); ok2 {
		userSessions := buf.([]DataChangeTrackerInterface)
		newSessions := []DataChangeTrackerInterface{}
		for _, session := range userSessions {
			if session.GetID() == us.GetID() {
				p.mylogger.Infof("remove session:%v for data type:%v", us.GetID(), dataType)
			} else {
				newSessions = append(newSessions, session)
			}
		}
		p.DataTopics.Store(dataType, newSessions)
		p.mylogger.Infof("data topics remove for a dataType:%v with len(UserSession):%v", dataType, len(newSessions))
	} else {
		p.mylogger.Warningf("skip to remove non-existing data topics:%v  session:%v", dataType, us.GetID())
	}
}

func (p *WebSocketManager) RemoveSession(sessionID string) {
	if sessBuf, ok := p.Sessions.Load(sessionID); ok {
		us, _ := sessBuf.(UserSessionInterface)
		if us != nil {
			p.removeSessionDataTopic(sessionID)
			us.Terminate()
		}
		p.Sessions.Delete(sessionID)
	} else {
		p.mylogger.Warningf("skip to remove session:%v because it is not load", sessionID)
	}
}

func (p *WebSocketManager) removeSessionDataTopic(sessionID string) {
	p.DataTopics.Range(func(dataType, value interface{}) bool {
		sessions := value.([]DataChangeTrackerInterface)
		newSessions := []DataChangeTrackerInterface{}
		for _, sess := range sessions {
			if sess.GetID() == sessionID {
				p.mylogger.Infof("try to delete session:%v", sess.GetID())
			} else {
				newSessions = append(newSessions, sess)
			}
		}
		p.DataTopics.Store(dataType, newSessions)
		return true
	})
}
