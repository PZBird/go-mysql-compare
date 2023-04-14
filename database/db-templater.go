package db

import (
	"fmt"

	"github.com/PZBird/go-mysql-compare/model"
)

func CreateDatabase(schema *model.DatabaseSchema) string {
	// CREATE SCHEMA `example`;
	query := fmt.Sprintf("CREATE SCHEMA `%s`;\n", schema.SchemaName)

	query += createTables(schema)

	return query
}

// CREATE TABLE `example`.`test` (
//
//	`id` int(10) unsigned NOT NULL AUTO_INCREMENT ,
//	`field1` varchar(45) NULL ,
//	`field2` tinyint(4) NULL ,
//	PRIMARY KEY (`id`)
//	);
func createTables(schema *model.DatabaseSchema) string {
	query := ""
	for _, table := range schema.Tables {
		query += fmt.Sprintf("CREATE TABLE `%s`.`%s` (\n", schema.SchemaName, table.TableName)
		for _, column := range table.Columns {
			query += addColumnForCreateTable(column)
		}

		query += addPrimaryColumn(table)

		query += ");\n"
	}

	return query
}

func addPrimaryColumn(table *model.Table) string {
	query := fmt.Sprintf("PRIMARY KEY (")

	needComma := false
	for _, column := range table.PrimaryKeys {
		query += fmt.Sprintf("`%s`", column.ColumnName)
		if needComma {
			query += ","
		}

		needComma = true
	}

	query += ")\n"
	return query
}

func addColumnForCreateTable(column *model.Column) string {
	query := fmt.Sprintf("`%s` %s", column.ColumnName, column.ColumnType)

	if !column.IsNullable {
		query += " NOT NULL"
	} else {
		query += " NULL"
	}

	if column.IsAutoIncrement {
		query += " AUTO_INCREMENT"
	}

	if column.DefaultValue != "" {
		query += fmt.Sprintf(" DEFAULT '%s'", column.DefaultValue)
	}

	query += ",\n"
	return query
}
