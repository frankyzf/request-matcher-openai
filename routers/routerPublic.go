package routers

import (
	"github.com/gin-gonic/gin"
	"request-matcher-openai/controls"
)

func routerPublic(router *gin.Engine) *gin.Engine {

	public := router.Group("/v1")
	{
		public.GET("/ws", controls.SubscribeHandler)
	}

	return router
}
