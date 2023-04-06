package main

import (
	"fmt"

	comparer "github.com/PZBird/go-mysql-compare/comparer"
	"github.com/PZBird/go-mysql-compare/configuration"
	db "github.com/PZBird/go-mysql-compare/database"
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

	if len(compareResult.LeftDatabaseExtraSchemas) > 0 {
		fmt.Println(fmt.Sprintf("New schemas for the %s:", hostnameLeft))

		for _, schema := range compareResult.LeftDatabaseExtraSchemas {
			query := db.CreateDatabase(schema)
			fmt.Println(query)
		}
	}

	if len(compareResult.RightDatabaseExtraSchemas) > 0 {
		fmt.Println(fmt.Sprintf("New schemas for the %s:", hostnameRight))

		for _, schema := range compareResult.RightDatabaseExtraSchemas {
			query := db.CreateDatabase(schema)
			fmt.Println(query)
		}
	}

	fmt.Println(compareResult)
}
