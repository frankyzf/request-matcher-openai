package routers

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) *gin.Engine {
	router = routerPublic(router)
	router = routerLogin(router)
	router = routerProject(router)

	return router
}
