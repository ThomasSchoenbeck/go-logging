package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	r.GET("/pingHTML", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ping.tmpl", gin.H{"title": "Ping", "response": "PONG!"})
	})

	r.POST("/logger", handleRemoteLogs)
	r.POST("/logs/:appID", handleLogsList)
	r.POST("/apps", handleAppsList)
	r.GET("/apps/:appID", handleGetAppByID)

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, token, X-Requested-With")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Content-Type", "Authorization", "Content-Length", "token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/stream", StreamHandler)

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	// r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 	// your custom format
	// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.ClientIP,
	// 		param.TimeStamp.Format(time.RFC1123),
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 		param.ErrorMessage,
	// 	)
	// }))
	r.Use(gin.Recovery())
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(fmt.Sprintf("0.0.0.0:%v", Port))
}
