package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
//	config.AllowOrigins = []string{"https://gin-template.vercel.app", "http://172.16.13.73:40080"}
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	return cors.New(config)
}
