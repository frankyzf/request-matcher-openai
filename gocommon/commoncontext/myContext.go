package commoncontext

import (
	"request-matcher-openai/gocommon/util"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"request-matcher-openai/gocommon/logger"
	"request-matcher-openai/gocommon/viperconf"
)

type MyContext struct {
	VP             *viper.Viper
	Name           string
	Loc            *time.Location
	LoggerLevel    string
	MyLogger       *logrus.Entry
	DBConn         *gorm.DB
	RClient        *redis.Client
	MyTokenManager *util.TokenManager
}

var instance *MyContext = nil

func GetInstance() *MyContext {
	return instance
}

func Setup(name string, vp *viper.Viper) error {
	var err error
	instance = &MyContext{}
	instance.VP = vp
	instance.Name = name
	instance.Loc, err = time.LoadLocation(viperconf.GetDefaultString(vp, "location", "Asia/Shanghai"))
	if err != nil {
		return err
	}
	instance.LoggerLevel = viperconf.GetDefaultString(vp, "logger_level", "info")
	instance.MyLogger = logger.SetupLogging(instance.LoggerLevel, "common", name)
	instance.DBConn, err = setupDatabaseConnection(vp)
	if err != nil {
		return err
	}
	instance.RClient, err = setupRedis(vp)
	if err != nil {
		return err
	}
	instance.MyTokenManager = util.GetTokenManager(instance.RClient, viperconf.GetDefaultInt(vp, "auth.token_expire_minute", 365*24*60))
	return nil
}

func Close() {
	if instance != nil {
		// instance.DBConn.Close()
	}
}

func GetLogger() *logrus.Entry {
	return instance.MyLogger
}
