package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// This function generates a connection string
//
// Parameters:
//
// user: string - username for database connection
//
// password: string - password for database connection
//
// host: string - host for database connection
//
// port: string - port for database connection
//
// Return string example:
// "root:password@tcp(127.0.0.1:3306)/"
func connectionString(connectionOptions ...interface{}) string {
	user := "root"
	password := "password"
	host := "127.0.0.1"
	port := "3306"

	for index, val := range connectionOptions {
		switch index {
		case 0: // check the user parameter
			user, _ = val.(string)
		case 1: // check the password parameter
			password, _ = val.(string)
		case 2: // check the host parameter
			host, _ = val.(string)
		case 3: // check the port parameter
			port, _ = val.(string)
		}
	}

	return user + ":" + password + "@tcp(" + host + ":" + port + ")/"
}

func connect(connectionString string) *sql.DB {
	log.SetPrefix("connection log: ")
	log.SetFlags(0)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	log.Println("Connected")
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
