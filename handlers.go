package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func handleLogsList(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeoutDuration)
	defer cancel()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")

	var lpr logsPaginationRequest
	if err := c.ShouldBindJSON(&lpr); err != nil {
		respondWithJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := getLogsPaginated(ctx, lpr)
	if err != nil {
		log.Println("error reading logs", err)
		respondWithJSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	respondWithJSON(c, http.StatusOK, res)

}

func handleRemoteLogs(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeoutDuration)
	defer cancel()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")

	// log.Println("client IP", c.ClientIP())
	// log.Println("fullpath", c.FullPath())
	// log.Println("remoteIP", c.RemoteIP())
	// log.Println("params", c.Params)
	// log.Println("param", c.Param(""))

	var logs []ClientLogs
	if err := c.ShouldBindJSON(&logs); err != nil {
		respondWithJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(logs); i++ {
		logs[i].CLIENT_IP = c.ClientIP()
		logs[i].REMOTE_IP = c.RemoteIP()
		logs[i].USERAGENT = c.Request.UserAgent()

	}

	if len(logs) > 0 {
		err := saveLogMessages(ctx, logs)
		if err != nil {
			errMsg := fmt.Sprintln("error saving log Messages", err)
			log.Println(errMsg)
			respondWithJSON(c, http.StatusInternalServerError, gin.H{"error": errMsg})
			// c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}
	}

	respondWithJSON(c, http.StatusOK, logs)
	// respondWithJSON(c, http.StatusOK, gin.H{"status": "done"})
	// c.JSON(http.StatusOK, gin.H{"status": "done"})
}

func respondWithJSON(c *gin.Context, statusCode int, message interface{}) {
	if strings.Contains(fmt.Sprint(message), "context deadline exceeded") {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": message})
	} else {
		c.JSON(statusCode, message)
	}
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

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()


	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" { {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG, PNG or GIF image", http.StatusBadRequest)
		return
	}

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")

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
