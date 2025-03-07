package main

import (
	"gin-template/common"
	"gin-template/middleware"
	"gin-template/model"
	"gin-template/router"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	common.SetupGinLog()
	common.SysLog("Gin Template " + common.Version + " started")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Initialize SQL Database
	err := model.InitDB()
	if err != nil {
		common.FatalLog(err)
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.FatalLog(err)
		}
	}()

	go model.Watch()

	// Initialize Redis
	err = common.InitRedisClient()
	if err != nil {
		common.FatalLog(err)
	}

	// Initialize options
	model.InitOptionMap()

	// Initialize HTTP server
	server := gin.Default()
	//server.Use(gzip.Gzip(gzip.DefaultCompression))
	server.Use(middleware.CORS())

	// Initialize session store
	if common.RedisEnabled {
		opt := common.ParseRedisOption()
		store, _ := redis.NewStore(opt.MinIdleConns, opt.Network, opt.Addr, opt.Password, []byte(common.SessionSecret))
		server.Use(sessions.Sessions("session", store))
	} else {
		store := cookie.NewStore([]byte(common.SessionSecret))
		store.Options(sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7, // Adjust based on your requirements
			Secure:   false,     // Ensure cookies are sent over HTTPS
			HttpOnly: true,      // Protect cookies from being accessed through client-side scripts
		})
		server.Use(sessions.Sessions("session", store))
	}

	router.SetRouter(server)
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}
	err = server.Run(":" + port)
	if err != nil {
		log.Println(err)
	}

}
