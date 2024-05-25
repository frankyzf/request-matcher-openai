package commoncontext

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"request-matcher-openai/gocommon/logger"
)

func StartLogLevelListener(vp *viper.Viper) {
	//signal handle
	name := GetInstance().Name
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGQUIT)
	go func() {
		for true {
			s := <-signalChannel
			GetInstance().MyLogger.Infof("recv a signal:%v", s)
			cf, _ := os.OpenFile(fmt.Sprintf("/tmp/%s.signal", name), os.O_RDWR|os.O_CREATE, 0644)
			defer cf.Close()
			buf, _ := ioutil.ReadAll(cf)
			command := string(buf)
			command = strings.TrimSpace(command)
			GetInstance().MyLogger.Infof("recv a signal command: %s", command)
			if command == "debug" {
				logger.SetLevel("debug")
			} else if command == "info" {
				logger.SetLevel("info")
			} else {
				GetInstance().MyLogger.Infof("unknown signal command:%v", command)
			}
		}
	}()
}
