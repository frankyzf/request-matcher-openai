package routers

import (
	"github.com/gin-gonic/gin"
	"request-matcher-openai/gocommon/jwt"

	"request-matcher-openai/controls"
	"request-matcher-openai/gocommon/commoncontext"
)

func routerLogin(router *gin.Engine) *gin.Engine {
	login := router.Group("/v1")
	{
		login.POST("/login", controls.LoginAccount)
		login.POST("/signup", controls.Signup)

		login.POST("/reset-password-with-email", controls.SignupUser)
		login.POST("/app-signup", controls.SignupUser)
		login.POST("/app-login", controls.LoginUser)
	}

	logout := router.Group("/v1")
	logout.Use(jwt.JWTAuth(commoncontext.GetDefaultString("cognito.secret", "mysecret")))
	{
		logout.POST("/logout", controls.Logout)
		logout.POST("/app-logout", controls.LogoutUser)
	}

	me := router.Group("/v1")
	me.Use(jwt.JWTAuth(commoncontext.GetDefaultString("cognito.secret", "mysecret")))
	{
		me.GET("/me", controls.GetMe)
		me.POST("/me/update-password", controls.UpdateMyPassword)
		me.POST("/me/self-delete-with-phone", controls.SelfDeleteWithPhone)
		me.POST("/me/self-delete-with-email", controls.SelfDeleteWithEmail)

		me.GET("/my-qualified-project-list", controls.GetMyQualifiedProjectList)
	}

	return router
}
