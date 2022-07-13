package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Port                   int
	Timeout                int
	ContextTimeoutDuration time.Duration
)

func processFlags() {
	log.Println("Process flags")
	flag.IntVar(&Port, "port", GetEnvAsInt("GOAPI_PORT", 8080), "Set API Port")
	flag.IntVar(&Port, "p", GetEnvAsInt("GOAPI_PORT", 8080), "Set API Port")
	flag.IntVar(&Timeout, "timeout", GetEnvAsInt("GOAPI_TIMEOUT", 10), "Set API Timeout")
	flag.IntVar(&Timeout, "t", GetEnvAsInt("GOAPI_TIMEOUT", 10), "Set API Timeout")
	flag.Parse()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	processFlags()
	checkDbConnection()
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		sig := <-c
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		db.Close()
		time.Sleep(1 * time.Second)
		// Alles aufraumen, commections schließen
		os.Exit(0)
	}()

	// Set Duration Time of the Context with timeout for api handlers
	ContextTimeoutDuration = time.Duration(Timeout) * time.Second

	setupRouter()

}
