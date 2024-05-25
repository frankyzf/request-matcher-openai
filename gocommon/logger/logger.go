package logger

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var LogInstance *log.Logger

func init() {
	if LogInstance != nil {
		return
	}
	LogInstance = log.New()

	path := path.Join("log/log-info")
	infoWriter, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName("log-info"),
		rotatelogs.WithMaxAge(90*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		fmt.Printf("can not open info logging file %v\n", err)
		os.Exit(1)
	}

	pathMap := lfshook.WriterMap{
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  infoWriter,
		log.ErrorLevel: infoWriter,
		log.DebugLevel: infoWriter,
	}

	LogInstance.Hooks.Add(lfshook.NewHook(
		pathMap,
		&log.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02T15:04:05.00+08:00",
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "@timestamp",
				log.FieldKeyLevel: "@level",
				log.FieldKeyMsg:   "@message",
			},
		},
	))

	if os.Getenv("env") == "prod" {
		SetLevel("info")
	} else {
		SetLevel("debug")
	}
}

func SetupLogging(env, owner, id string) *log.Entry {
	if env == "" {
		env = "unkonwn"
	}
	SetLevel(env)
	// contextLogger := LogInstance.WithFields(log.Fields{
	// 	"env":   env,
	// 	"owner": owner,
	// 	"id":    id,
	// })
	contextLogger := LogInstance.WithFields(log.Fields{
		"owner": owner,
		"id":    id,
	})
	return contextLogger
}

func SetLevel(level string) {
	lv := log.InfoLevel
	if strings.ToUpper(level) == "DEBUG" {
		lv = log.DebugLevel
	} else {
		lv = log.InfoLevel
	}
	LogInstance.SetLevel(lv)
}

func IsDebugLevel() bool {
	lv := LogInstance.GetLevel()
	if lv == log.DebugLevel {
		return true
	}
	return false
}
