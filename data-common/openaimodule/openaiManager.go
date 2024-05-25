package openaimodule

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"request-matcher-openai/gocommon/commoncontext"
)

type OpenaiManager struct {
	openaiAdapter        *OpenaiAdapter
	myRClient            *redis.Client
	mylogger             *logrus.Entry
	appID                string
	secret               string
	appTokenExpireMinute int64
	loggerLevel          string
}

func GetOpenaiManager(rclient *redis.Client) *OpenaiManager {
	p := &OpenaiManager{
		myRClient: rclient,
		mylogger:  commoncontext.SetupLogging("openai_module", "openai_manager"),
	}

	accessTokenURL := commoncontext.GetDefaultString("openaimodule.access_token_url", "")
	liveDataBaseURL := commoncontext.GetDefaultString("openaimodule.live_data_base_url", "")
	p.openaiAdapter = GetOpenaiAdapter(accessTokenURL, liveDataBaseURL)

	p.appID = commoncontext.GetDefaultString("openaimodule.app_id", "")
	p.secret = commoncontext.GetDefaultString("openaimodule.secret", "")
	p.appTokenExpireMinute = int64(commoncontext.GetDefaultInt("openaimodule.token_expire_minute", 120))
	return p
}

func (p *OpenaiManager) GetName() string {
	return "openai_manager"
}
