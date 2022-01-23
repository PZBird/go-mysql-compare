package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/PZBird/go-mysql-compair/configuration"
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
func ConnectionString(connectionOptions *configuration.ConfigurationElement) string {
	connectionString := ""
	connectionString += connectionOptions.Username
	connectionString += ":"
	connectionString += connectionOptions.Password
	connectionString += "@tcp("
	connectionString += connectionOptions.Hostname
	connectionString += ":"
	connectionString += connectionOptions.Port
	connectionString += ")/information_schema?parseTime=true"

	return connectionString
}

// Connect to the database using a connection string. Like: "username:password@tcp(address:port)/"
func Connect(connectionString string) *sql.DB {
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
