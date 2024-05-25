package util

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func GetMyGinRouter() *gin.Engine {
	router := gin.Default()
	ReqHeader := []string{
		"Content-Type", "Origin", "Authorization", "Accept", "token",
		"cache-control", "x-requested-with", "x-token", "lat", "lon"}
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, PATCH",
		RequestHeaders:  strings.Join(ReqHeader, ", "),
		ExposedHeaders:  "",
		MaxAge:          24 * 3600 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	return router
}
