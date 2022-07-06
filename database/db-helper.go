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

		databaseWithoutSuffix := strings.ReplaceAll(databaseName, databasesSuffix, "")

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

type indexResponse struct {
	IndexName  string
	SeqInIndex int8
	ColumnName string
	NonUnique  bool
	IndexType  string
	Comment    string
}

func readIndexes(conn *sql.DB, schema *model.DatabaseSchema,
	tableName string, table *model.Table, hostname string) {
	q := "SELECT"
	q += "	index_name,"
	q += "	seq_in_index,"
	q += "	column_name,"
	q += "	non_unique,"
	q += "	index_type,"
	q += "	comment"
	q += " FROM"
	q += "		INFORMATION_SCHEMA.STATISTICS"
	q += " WHERE 1=1"
	q += "		AND table_schema = ?"
	q += "		AND table_name = ?"
	q += " ORDER BY seq_in_index;"

	rows, err := conn.Query(q, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		index := &model.Index{}
		index.TableName = tableName
		var response indexResponse

		err := rows.Scan(
			&response.IndexName,
			&response.SeqInIndex,
			&response.ColumnName,
			&response.NonUnique,
			&response.IndexType,
			&response.Comment,
		)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func readColumns(conn *sql.DB, schema *model.DatabaseSchema,
	tableName string, table *model.Table, hostname string) {

	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, "
	q += " CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, "
	q += " COLUMN_TYPE, COLUMN_KEY, EXTRA"
	q += " FROM information_schema.COLUMNS "
	q += "WHERE TABLE_SCHEMA=? AND TABLE_NAME=? ORDER BY ORDINAL_POSITION;"

	rows, err := conn.Query(q, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		column := &model.Column{}
		column.DatabaseName = hostname
		nullable := "NO"
		columnKey, extra := "", ""
		err := rows.Scan(&column.TableName, &column.ColumnName,
			&nullable, &column.DataType,
			&column.CharacterMaximumLength, &column.NumericPrecision,
			&column.NumericScale, &column.ColumnType, &columnKey, &extra)
		if err != nil {
			log.Fatal(err)
		}

		table.Columns = make(map[string]*model.Column)

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

		columnName := column.ColumnName

		table.Columns[columnName] = column
	}
}
