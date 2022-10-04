package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db        *sql.DB
	tableList = [...]string{"TSC_APPLICATIONS", "TSC_CLIENT_LOGS", "TSC_FEEDBACK_CHANNELS", "TSC_FEEDBACK"}
)

// func checkDbConnection() {
// 	dbUser := os.Getenv("DB_USER")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dsn := fmt.Sprintf("hdb://%s:%s@%s:%s", dbUser, url.QueryEscape(os.Getenv("DB_PASSWORD")), dbHost, dbPort)
// 	log.Println("Checking DB connection with: " + dbHost + ":" + dbPort)

// 	var errOpen error
// 	db, errOpen = sql.Open("hdb", dsn)
// 	if errOpen != nil {
// 		log.Fatal(errOpen)
// 	}

// 	errPing := db.Ping()
// 	if errPing != nil {
// 		log.Fatal(errPing)
// 	}
// 	log.Println("Connection successful with DB: " + dbHost + ":" + dbPort)
// 	log.Printf("hdb://%s:%s@%s:%s\n", dbUser, "**************", dbHost, dbPort)
// }

func checkDbConnection() {
	var err error
	db, err = sql.Open("sqlite3", "foo.db")
	if err != nil {
		log.Println("error opening sqlite db")
	}
	log.Println("DB Connection successful")
}

func checkDbIntegrity() {

	// ctx, cancel := context.WithTimeout(context.Background(), ContextTimeoutDuration)
	// defer cancel()
	// always error contect deadline exeeded when timeout = 10 seconds
	fileName := "create_tables.sql"
	dat, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("error reading file", fileName)
	}

	sqlStatements := strings.Split(string(dat), ";")

	var statements []string
	for _, s := range sqlStatements {
		if len(s) > 0 {
			statements = append(statements, s)
		}
	}

	log.Println("statements", len(statements))

	for _, t := range tableList {
		var numberOfRecords int
		err := db.QueryRow(fmt.Sprintf("SELECT count(*) from %s", t)).Scan(&numberOfRecords)
		if err != nil {
			log.Println("Table", t, "is missing: ", err)
			for _, s := range statements {
				if strings.Contains(s, t+" (") {
					_, err := db.Exec(s)
					if err != nil {
						log.Println("error creating missing table", t, err)
					}
					log.Println("created missing table", t)
				}
			}
		}
	}
}
