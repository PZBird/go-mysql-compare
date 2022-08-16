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

	hostname1 := config.Db1.Hostname + ":" + config.Db1.Port
	hostname2 := config.Db2.Hostname + ":" + config.Db2.Port

	databaseTablesFromDb1 := db.GetDatabaseTablesOrFail(db1, config.Db1.DatabasesSuffix, hostname1)
	databaseTablesFromDb2 := db.GetDatabaseTablesOrFail(db2, config.Db2.DatabasesSuffix, hostname2)

	compareResult := comparer.Compare(databaseTablesFromDb1, databaseTablesFromDb2, config)

	fmt.Println(compareResult)
}
