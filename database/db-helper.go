package db

import (
	"database/sql"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/PZBird/go-mysql-compair/model"
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
func GetDatabaseTablesOrFail(db *sql.DB) map[string]*model.DatabaseSchema {
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

		schema := &model.DatabaseSchema{}
		schema.SchemaName = databaseName

		readTables(db, schema)

		databaseTables[databaseName] = schema
	}

	return databaseTables
}

func readTables(conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.TABLES "
	q += "Where TABLE_SCHEMA=?"
	q += " AND TABLE_TYPE='BASE TABLE' ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	schema.Tables = make(map[string]*model.Table)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		table := &model.Table{}
		err := rows.Scan(&table.TableName)

		if err != nil {
			log.Fatal(err)
		}

		schema.Tables[table.TableName] = table

		log.Printf("Examining table %s\r\n", table.TableName)
		readColumns(conn, schema, table.TableName, &table.Columns)
		for _, col := range table.Columns {
			if col.IsPrimaryKey {
				table.PrimaryKeys = append(table.PrimaryKeys, col)
			} else {
				table.OtherColumns = append(table.OtherColumns, col)
			}
		}
	}
}

func readColumns(conn *sql.DB, schema *model.DatabaseSchema,
	tableName string, columns *[]*model.Column) {
	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, "
	q += " CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, "
	q += " COLUMN_TYPE, COLUMN_KEY, EXTRA"
	q += " FROM information_schema.COLUMNS "
	q += "WHERE TABLE_SCHEMA=? AND TABLE_NAME=? ORDER BY ORDINAL_POSITION"
	rows, err := conn.Query(q, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		column := &model.Column{}
		nullable := "NO"
		columnKey, extra := "", ""
		err := rows.Scan(&column.TableName, &column.ColumnName,
			&nullable, &column.DataType,
			&column.CharacterMaximumLength, &column.NumericPrecision,
			&column.NumericScale, &column.ColumnType, &columnKey, &extra)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Examining column %s\r\n", column.ColumnName)

		if column.DataType == "enum" {
			regInBrackets := regexp.MustCompile(`\((.*?)\)`)
			enumValues := regInBrackets.FindStringSubmatch(column.ColumnType)
			values := strings.Split(enumValues[1], ",")
			sort.Strings(values)
			column.EnumValues = values
		}

		if nullable == "NO" {
			column.IsNullable = false
		} else {
			column.IsNullable = true
		}

		if columnKey == "PRI" {
			column.IsPrimaryKey = true
		}

		if columnKey == "UNI" {
			column.IsUnique = true
		}

		if extra == "auto_increment" {
			column.IsAutoIncrement = true
		}

		*columns = append(*columns, column)
	}
}
