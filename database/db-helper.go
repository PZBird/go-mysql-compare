package db

import (
	"database/sql"
	"log"
)

// LIST OF EXCLUDED DATABASES
var EXCLUDE_LIST = map[string]bool{
	"information_schema": true,
	"mysql":              true,
	"performance_schema": true,
	"sys":                true,
}

func databaseInBlackList(databaseName string) bool {
	_, exists := EXCLUDE_LIST[databaseName]

	return exists
}

func GetDatabasesOrFail(db *sql.DB) []string {
	var databaseName string
	var databaseNames []string

	rows, error := db.Query("SHOW DATABASES")

	if error != nil {
		panic(error)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&databaseName)

		if err != nil {
			log.Fatal(err)
		}

		databaseNames = append(databaseNames, databaseName)
	}

	return databaseNames
}

// Return a map where key is database name and value is array of table names
func GetDatabaseTablesOrFail(db *sql.DB) map[string][]string {
	var table string
	var database string
	databaseTables := make(map[string][]string)

	q := "SELECT TABLE_SCHEMA, TABLE_NAME FROM information_schema.TABLES"
	q += " WHERE TABLE_TYPE='BASE TABLE'"
	q += " AND TABLE_SCHEMA NOT IN ("
	q += "'information_schema','mysql','performance_schema','sys'"
	q += ") ORDER BY TABLE_NAME"

	rows, error := db.Query(q)

	if error != nil {
		panic(error)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&database, &table)

		if err != nil {
			log.Fatal(err)
		}

		databaseTables[database] = append(databaseTables[database], table)
	}

	return databaseTables
}
