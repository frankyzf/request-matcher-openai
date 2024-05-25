package main

import (
	_ "net/http/pprof"
	"os"
	"request-matcher-openai/data-common/mywebsocket"
	"request-matcher-openai/data-mydb/export"
	"time"

	"request-matcher-openai/controls"
	"request-matcher-openai/data-common/datamanage"
	"request-matcher-openai/data-common/myauth"
	"request-matcher-openai/data-mydb/dmmodule"
	"request-matcher-openai/gocommon/util"
	"request-matcher-openai/routers"

	_ "github.com/lib/pq" // here
	"github.com/sirupsen/logrus"

	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/viperconf"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var vp *viper.Viper
var auth *myauth.MyAuth
var mylogger *logrus.Entry

func main() {
	var err error
	var Env string
	{
		pflag.StringVar(&Env, "env", "dev", "environment")
		pflag.String("dialect", "mysql", "mysql/postgres")
		pflag.String("location", "Asia/Singapore", "timezone locaiton")
		pflag.String("secret", "mysecret", "secret for authoriazation header")
		pflag.String("captchasecret", "6LexTe0UAAAAADnTGZ81QJID1Zrv3Y75F91DJ44A", "captcha string")
		pflag.Bool("debug", true, "turn on debug")
	}
	pflag.Parse()

	vp = viper.New()
	vp.BindPFlags(pflag.CommandLine)
	vp = commoncontext.GetAndMergeConfig(vp, Env)

	if viperconf.GetDefaultBool(vp, "gin_debug", true) == false {
		gin.SetMode(gin.ReleaseMode)
	}

	err = commoncontext.Setup("request-matcher-openai", vp)
	mylogger = commoncontext.GetInstance().MyLogger
	if err != nil {
		mylogger.Errorf("failed to setup common context:%v", err)
		os.Exit(1)
	}
	defer commoncontext.Close()
	mylogger.Infof("start process request-matcher-openai:%v", time.Now())

	dbconn := commoncontext.GetInstance().DBConn

	err = MyInit(dbconn)
	if err != nil {
		mylogger.Errorf("failed to init:%v", err)
		os.Exit(1)
	}

	auth = myauth.Setup()

	dmManage, err4 := dmmodule.Setup()
	if err4 != nil {
		mylogger.Errorf("failed to get dmmanage:%v", err4)
		os.Exit(1)
	}

	manage, err3 := datamanage.Setup(auth, dmManage)
	if err3 != nil {
		mylogger.Errorf("failed to get datamanage:%v", err)
		os.Exit(1)
	}

	//call the func after manage is ready
	dmManage.Init()
	// auth.SetupDMManage(dmManage)

	var wsManager export.WebSocketManagerInterface
	if viperconf.GetDefaultBool(vp, "websocket.enable", false) {
		wsManager, err = mywebsocket.Setup(manage, auth)
		if err != nil {
			mylogger.Errorf("failed to setup websocket, err:%v", err)
			os.Exit(1)
		}
		wsManager.Start()
		manage.SetupWebSocketManager(wsManager)
	} else {
		mylogger.Warningf("skip to setup websocket because websocket.enable is turned off")
	}

	controls.Setup(auth, manage, wsManager)

	router := util.GetMyGinRouter()
	router.Use(RequestLogger())

	router = routers.Register(router)

	commoncontext.StartLogLevelListener(vp)

	// buildMachine, _ := os.Hostname()
	mylogger.WithField("op", "start").Infof("start app version:%s", versionString())

	mylogger.Fatal(router.Run(vp.GetString("addr")))

}
