package main

import (
	"fmt"
)

func main() {
	db := connect(connectionString())

	databaseTables := getDatabaseTablesOrFail(db)

	fmt.Println(databaseTables)
}
