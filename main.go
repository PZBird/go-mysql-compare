package main

import (
	"fmt"

	"github.com/PZBird/go-mysql-compair/configuration"
	"github.com/PZBird/go-mysql-compair/database"
)

func main() {
	config := configuration.LoadConfiguration("./config.json")

	db1 := db.Connect(db.ConnectionString(&config.Db1))
	db2 := db.Connect(db.ConnectionString(&config.Db2))

	databaseTablesFromDb1 := db.GetDatabaseTablesOrFail(db1)
	databaseTablesFromDb2 := db.GetDatabaseTablesOrFail(db2)

	fmt.Println(databaseTablesFromDb1)
	fmt.Println(databaseTablesFromDb2)
}
