package main

import (
	"fmt"
	"strings"

	comparer "github.com/PZBird/go-mysql-compare/comparer"
	"github.com/PZBird/go-mysql-compare/configuration"
	db "github.com/PZBird/go-mysql-compare/database"
	"github.com/PZBird/go-mysql-compare/model"
)

func main() {
	config := configuration.LoadConfiguration("./config.json")

	db1 := db.Connect(db.ConnectionString(&config.Db1))
	db2 := db.Connect(db.ConnectionString(&config.Db2))

	hostnameLeft := config.Db1.Hostname + ":" + config.Db1.Port
	hostnameRight := config.Db2.Hostname + ":" + config.Db2.Port

	databaseLeft := db.GetDatabaseTablesOrFail(db1, config.Db1.DatabasesSuffix, hostnameLeft)
	databaseRight := db.GetDatabaseTablesOrFail(db2, config.Db2.DatabasesSuffix, hostnameRight)

	compareResult := comparer.Compare(databaseLeft, databaseRight, config)

	genCteateDatabaseScripts(compareResult.LeftDatabaseExtraSchemas, hostnameLeft)
	genCteateDatabaseScripts(compareResult.RightDatabaseExtraSchemas, hostnameRight)

	leftColumnsSql, rightColumsSql := genColumnDiffScripts(&compareResult)

	fmt.Println(leftColumnsSql)
	fmt.Println(rightColumsSql)

	fmt.Println(compareResult)
}

func genCteateDatabaseScripts(databaseExtraSchemas []*model.DatabaseSchema, hostnameLeft string) {
	if len(databaseExtraSchemas) > 0 {
		fmt.Printf("New schemas for the %s:\n", hostnameLeft)

		for _, schema := range databaseExtraSchemas {
			query := db.CreateDatabase(schema)
			fmt.Println(query)
		}
	}
}

func genColumnDiffScripts(result *comparer.ComparerResult) (leftColumSql, reightColumnSql string) {
	var leftColumnSqlBuilder strings.Builder
	var rightColumnSqlBuilder strings.Builder

	db.GenNewColumnsScripts(&leftColumnSqlBuilder, result.ColumnToInsertDBLeft)
	db.GenNewColumnsScripts(&rightColumnSqlBuilder, result.ColumnToInsertDBRight)

	db.GenRemovedColumnsScripts(&leftColumnSqlBuilder, result.ColumnToRemoveDBLeft)
	db.GenRemovedColumnsScripts(&rightColumnSqlBuilder, result.ColumnToRemoveDBRight)

	db.GenModifiedColumnsScripts(&leftColumnSqlBuilder, result.ColumnToModifyDBLeft)
	db.GenModifiedColumnsScripts(&rightColumnSqlBuilder, result.ColumnToModifyDBRight)

	return leftColumnSqlBuilder.String(), rightColumnSqlBuilder.String()
}
