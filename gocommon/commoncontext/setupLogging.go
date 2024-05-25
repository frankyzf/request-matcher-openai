package commoncontext

import (
	"github.com/sirupsen/logrus"
	"request-matcher-openai/gocommon/logger"
)

func SetupLogging(module, name string) *logrus.Entry {
	loggerLevel := GetInstance().LoggerLevel
	mylogger := logger.SetupLogging(loggerLevel, module, name)
	return mylogger
}
