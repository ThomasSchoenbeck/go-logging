package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	Port = 8080
)

func init() {

}

func randomNumber() int {
	min := 10
	max := 250

	return rand.Intn(max-min) + min
}

func StreamHandler(c *gin.Context) {
	randNum := randomNumber()
	c.String(http.StatusOK, fmt.Sprintf("random number: %d", randNum))
}

// func StreamHandler(w http.ResponseWriter, r *http.Request) {
// 	flusher, ok := w.(http.Flusher)

// 	if !ok {
// 		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	for i := 0; i < 20; i++ {
// 		randNum := randomNumber()
// 		fmt.Fprintf(w, "Index: %d  ->  waited for %dms\n", i, randNum)
// 		// msg := sse{Id: i, Method: http.StatusOK, Body: "body data", Data: fmt.Sprintf("Index: %d  ->  waited for %dms\n", i, randNum), Time: time.Now().String(), Type: "message"}
// 		// json.NewEncoder(w).Encode(msg)
// 		flusher.Flush()
// 		time.Sleep(time.Duration(randNum) * time.Millisecond)
// 	}
// 	// w.WriteHeader(http.StatusInternalServerError)

// 	fmt.Fprintln(w, "done")
// }

type remoteLog struct {
	Logs []logs `json:"logs"`
}

type logs struct {
	Msg        string    `json:"msg"`
	Level      string    `json:"level"`
	Stacktrace string    `json:"stacktrace"`
	Timestamp  time.Time `json:"timestamp"`
}

func handleRemoteLogs(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	// log.Println("c.FullPath", c.FullPath())

	// var json Login
	// if err := c.ShouldBindJSON(&json); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	var rl remoteLog
	if err := c.ShouldBindJSON(&rl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("remote log: %#v \n", rl)

	c.JSON(http.StatusOK, gin.H{"done": "done"})
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	r.GET("/pingHTML", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ping.tmpl", gin.H{"title": "Ping", "response": "PONG!"})
	})

	r.POST("/logger", handleRemoteLogs)

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, token, X-Requested-With")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Content-Length", "token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
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
