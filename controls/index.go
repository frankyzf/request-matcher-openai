package controls

import (
	"errors"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"

	"request-matcher-openai/data-common/datamanage"
	"request-matcher-openai/data-common/myauth"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/util"
)

var mylogger *log.Entry

var myRClient *redis.Client
var auth *myauth.MyAuth
var manage *datamanage.DataManager
var myWebSocketManager export.WebSocketManagerInterface
var myTokenManager *util.TokenManager

func Setup(mauth *myauth.MyAuth, dm *datamanage.DataManager, ws export.WebSocketManagerInterface) error {
	mylogger = commoncontext.SetupLogging("request-matcher-openai", "controls")
	myRClient = commoncontext.GetInstance().RClient
	myTokenManager = commoncontext.GetInstance().MyTokenManager
	auth = mauth
	manage = dm
	myWebSocketManager = ws

	return nil
}

func SetupLoginToken(userID string, accountType string, expireStr string, token string) error {
	if expireStr != "" {
		expire, err := strconv.Atoi(expireStr)
		if err != nil {
			return err
		}
		myTokenManager.SetLoginToken(userID, accountType, int64(expire), token)
		return err
	} else {
		return errors.New("empty expire second")
	}
}
