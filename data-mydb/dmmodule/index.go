package dmmodule

import (
	"runtime"
	"strings"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"request-matcher-openai/gocommon/commoncontext"
)

var loc *time.Location
var myDbConn *gorm.DB
var myRClient *redis.Client
var mylogger *log.Entry

// var tableNotifySender export.NotifySenderInterface

func Setup() (*DMManage, error) {
	mylogger = commoncontext.SetupLogging("dmmodule", "manage")
	myDbConn = commoncontext.GetInstance().DBConn
	myRClient = commoncontext.GetInstance().RClient
	// tableNotifySender = tns
	loc = commoncontext.GetInstance().Loc

	// notifierListenKey := commoncontext.GetDefaultString("cache.notifier_key", "notifier_handler")
	p := GetDMManage(myDbConn, myRClient)
	return p, nil
}

func GetCurrentFunctionName() string {
	// Skip GetCurrentFunctionName
	fn := getFrame(1).Function
	index := strings.LastIndex(fn, "/")
	if index > 0 {
		index++
	}
	return fn[index:len(fn)]
}

func GetCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	fn := getFrame(2).Function
	index := strings.LastIndex(fn, "/")
	if index > 0 {
		index++
	}
	return fn[index:len(fn)]
}

func GetParentCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	fn := getFrame(3).Function
	index := strings.LastIndex(fn, "/")
	if index < 0 {
		return "ROOT"
	}
	if index > 0 {
		index++
	}
	return fn[index:len(fn)]
}

func GetGrandCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	fn := getFrame(4).Function
	index := strings.LastIndex(fn, "/")
	if index < 0 {
		return "ROOT"
	}
	if index > 0 {
		index++
	}
	return fn[index:len(fn)]
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}
	return frame
}
