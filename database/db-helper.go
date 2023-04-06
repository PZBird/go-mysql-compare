package db

import (
	"database/sql"
	"log"
	"strings"

	"github.com/PZBird/go-mysql-compare/model"
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
func GetDatabaseTablesOrFail(db *sql.DB, databasesSuffix string, hostname string) map[string]*model.DatabaseSchema {
	var tableName string
	var databaseName string
	databaseTables := make(map[string]*model.DatabaseSchema)

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
		err := rows.Scan(&databaseName, &tableName)

		if err != nil {
			log.Fatal(err)
		}

		databaseWithoutSuffix := strings.TrimSuffix(databaseName, databasesSuffix)

		schema := &model.DatabaseSchema{}
		schema.SchemaName = databaseName
		schema.DatabaseName = hostname

		readTables(db, schema, hostname)

		databaseTables[databaseWithoutSuffix] = schema
	}

	return databaseTables
}

func readTables(conn *sql.DB, schema *model.DatabaseSchema, hostname string) {
	q := "SELECT TABLE_NAME FROM information_schema.TABLES "
	q += "WHERE TABLE_SCHEMA=?"
	q += " AND TABLE_TYPE='BASE TABLE' ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	schema.Tables = make(map[string]*model.Table)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		table := &model.Table{}
		table.DatabaseName = hostname
		err := rows.Scan(&table.TableName)

		if err != nil {
			log.Fatal(err)
		}

		schema.Tables[table.TableName] = table

		readColumns(conn, schema, table.TableName, table, hostname)
		readIndexes(conn, schema, table.TableName, table, hostname)

		for _, col := range table.Columns {
			if col.IsPrimaryKey {
				table.PrimaryKeys = append(table.PrimaryKeys, col)
			} else {
				table.OtherColumns = append(table.OtherColumns, col)
			}
		}
	}
}
