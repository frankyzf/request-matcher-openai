package datamanage

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"request-matcher-openai/data-common/myauth"
	"request-matcher-openai/data-mydb/dmmodule"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/commoncontext"
)

var loc *time.Location
var myDbConn *gorm.DB
var myAuth *myauth.MyAuth
var myRClient *redis.Client
var mylogger *log.Entry

type DataManager struct {
	DMManage           *dmmodule.DMManage
	myWebSocketManager export.WebSocketManagerInterface //init laster by calling SetupWebSocketManager
	myAuth             *myauth.MyAuth
	myDbConn           *gorm.DB
	myRClient          *redis.Client
	mylogger           *logrus.Entry
	loggerLevel        string
}

func Setup(mauth *myauth.MyAuth, dmm *dmmodule.DMManage) (*DataManager, error) {
	mylogger = commoncontext.SetupLogging("0", "manage")
	myDbConn = commoncontext.GetInstance().DBConn
	myRClient = commoncontext.GetInstance().RClient
	myAuth = mauth
	loc = commoncontext.GetInstance().Loc

	p := &DataManager{
		DMManage:  dmm,
		myAuth:    mauth,
		myDbConn:  myDbConn,
		myRClient: myRClient,
		mylogger:  mylogger,
	}

	return p, nil
}

func (p *DataManager) GetName() string {
	return "datamanager"
}

func (p *DataManager) SetupWebSocketManager(ws export.WebSocketManagerInterface) {
	p.mylogger.Infof("websocket is setup in data maanger")
	p.myWebSocketManager = ws
}

func (p *DataManager) GetDBConn() *gorm.DB {
	return p.myDbConn
}

func (p *DataManager) GetRawCount(sql string) int {
	count := struct {
		Count int64
	}{}
	p.myDbConn.Raw(sql).Scan(&count)
	return int(count.Count)
}

func (p *DataManager) GenerateRandomCode() string {
	num := rand.Intn(1000000)
	code := fmt.Sprintf("%06d", num)
	env := commoncontext.GetDefaultString("env", "uat")
	if env == "uat" || env == "stage" || env == "dev" {
		p.mylogger.Infof("mock env, verification code is changed from :%v to :%v", code, "666666")
		code = "666666"
	}
	return code
}
