package openaimodule

import (
	"request-matcher-openai/gocommon/commoncontext"
)

var myOpenaiManager *OpenaiManager

func Setup() *OpenaiManager {
	if myOpenaiManager == nil {
		myOpenaiManager = GetOpenaiManager(commoncontext.GetInstance().RClient)
	}

	return myOpenaiManager
}
