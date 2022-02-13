package comparer

import (
	"fmt"
	"log"

	"github.com/PZBird/go-mysql-compair/configuration"
	"github.com/PZBird/go-mysql-compair/model"
)

type ComparerResult struct {
	LeftDatabaseExtraSchemas  []*model.DatabaseSchema
	RightDatabaseExtraSchemas []*model.DatabaseSchema
	TablesToInsertDB2         []*model.Table
}

func CompareSchemas(databaseTablesFromDb1 map[string]*model.DatabaseSchema, databaseTablesFromDb2 map[string]*model.DatabaseSchema, configuration configuration.Configuration) {
	var comparerResult ComparerResult

	for k, schemaFromDb1 := range databaseTablesFromDb1 {
		schemaFromDb2, isExist := databaseTablesFromDb2[k]

		// DB2 doesn't have table
		if !(isExist) {
			log.Print(fmt.Sprintf("Schema %s doesn't exist in compared db.", schemaFromDb1.SchemaName))
			comparerResult.LeftDatabaseExtraSchemas = append(comparerResult.LeftDatabaseExtraSchemas, schemaFromDb1)

			continue
		}

		compareTables(schemaFromDb1.Tables, schemaFromDb2.Tables, &comparerResult)
	}

	fmt.Println(comparerResult)
}

func compareTables(tables1 map[string]*model.Table, tables2 map[string]*model.Table, comparerResult *ComparerResult) {
	for tableName, tableStruct := range tables1 {
		tableStructFromDb2, isExist := tables2[tableName]

		if !(isExist) {
			log.Print(fmt.Sprintf("Table %s doesn't exist in compared db.", tableStruct.TableName))
			comparerResult.TablesToInsertDB2 = append(comparerResult.TablesToInsertDB2, tableStruct)

			continue
		}

		fmt.Println(tableStructFromDb2)
	}
}
