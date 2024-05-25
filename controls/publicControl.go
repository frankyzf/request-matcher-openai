package controls

import (
	"errors"

	"github.com/gin-gonic/gin"
	"request-matcher-openai/gocommon/replyutil"
)

// @Summary SubscribeHandler Info
// @Description SubscribeHandler
// @ID SubscribeHandler
// @Accept  json
// @Produce  json
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Router /ws [get]
func SubscribeHandler(c *gin.Context) {
	var err error
	if myWebSocketManager != nil {
		err = myWebSocketManager.AcceptSession(c)
	} else {
		err = errors.New("websocket server is not enabled")
	}

	if err != nil {
		replyutil.ResErr(c, err)
	}
}
