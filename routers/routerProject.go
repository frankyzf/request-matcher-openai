package routers

import (
	"github.com/gin-gonic/gin"

	"request-matcher-openai/controls"
	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/jwt"
)

func routerProject(router *gin.Engine) *gin.Engine {
	project := router.Group("/v1/project")
	project.Use(jwt.JWTAuth(commoncontext.GetDefaultString("cognito.secret", "mysecret")))
	{
		project.GET("/list", controls.GetProjectList)
		project.GET("/item/:id", controls.GetProjectItem)
		project.POST("/item", controls.CreateProjectItem)
		project.POST("/item/:id", controls.UpdateProjectItem)
		project.DELETE("/item/:id", controls.DeleteProjectItem)
	}

	return router
}
