package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
)

var (
	db *sql.DB
)

func checkDbConnection() {
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("hdb://%s:%s@%s:%s", dbUser, url.QueryEscape(os.Getenv("DB_PASSWORD")), dbHost, dbPort)
	log.Println("Checking DB connection with: " + dbHost + ":" + dbPort)

	var errOpen error
	db, errOpen = sql.Open("hdb", dsn)
	if errOpen != nil {
		log.Fatal(errOpen)
	}

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal(errPing)
	}
	log.Println("Connection successful with DB: " + dbHost + ":" + dbPort)
	log.Printf("hdb://%s:%s@%s:%s\n", dbUser, "**************", dbHost, dbPort)
}
