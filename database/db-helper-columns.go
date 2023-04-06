package db

import (
	"database/sql"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/PZBird/go-mysql-compare/model"
)

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

	table.Columns = make(map[string]*model.Column)

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

		appendColumn(column, nullable, columnKey, extra, table)
	}
}

func appendColumn(column *model.Column, nullable string, columnKey string, extra string, table *model.Table) {
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
