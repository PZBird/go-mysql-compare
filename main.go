package main

import (
	comparer "github.com/PZBird/go-mysql-compair/comparer"
	"github.com/PZBird/go-mysql-compair/configuration"
	db "github.com/PZBird/go-mysql-compair/database"
)

func main() {
	config := configuration.LoadConfiguration("./config.json")

	db1 := db.Connect(db.ConnectionString(&config.Db1))
	db2 := db.Connect(db.ConnectionString(&config.Db2))

	databaseTablesFromDb1 := db.GetDatabaseTablesOrFail(db1, config.Db1.DatabasesSuffix)
	databaseTablesFromDb2 := db.GetDatabaseTablesOrFail(db2, config.Db2.DatabasesSuffix)

	comparer.CompareSchemas(databaseTablesFromDb1, databaseTablesFromDb2, config)
}
