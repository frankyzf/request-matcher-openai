package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"request-matcher-openai/gocommon/replyutil"
)

func StaticFileAuth(secret string, tokenVerify bool, nextHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if tokenVerify {
			_, err := VerifyAuthorization(token, secret)
			if err != nil {
				fmt.Printf("static auth token err:%v\n", err)
				err = replyutil.StaticFileAuthFailure{Message: ""}
				replyutil.ResAppErr(c, err)
				c.AbortWithStatus(200)
				return
			}
		}
		nextHandler(c)
	}
}
