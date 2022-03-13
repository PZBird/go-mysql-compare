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
	TablesToInsertDBLeft      []*model.Table
	TablesToInsertDBRight     []*model.Table
	ColumnToInsertDBLeft      []*model.Column
	ColumnToInsertDBRight     []*model.Column
}

func Compare(databaseTablesFromDb1 map[string]*model.DatabaseSchema, databaseTablesFromDb2 map[string]*model.DatabaseSchema, configuration configuration.Configuration) ComparerResult {
	var comparerResult ComparerResult

	compareSchemas(databaseTablesFromDb1, databaseTablesFromDb2, configuration, &comparerResult, true)
	compareSchemas(databaseTablesFromDb2, databaseTablesFromDb1, configuration, &comparerResult, false)

	return comparerResult
}

func compareSchemas(databaseTablesFromDbLeft map[string]*model.DatabaseSchema, databaseTablesFromDbRight map[string]*model.DatabaseSchema, configuration configuration.Configuration, comparerResult *ComparerResult, isSideLeft bool) {
	for k, schemaFromDbLeft := range databaseTablesFromDbLeft {
		schemaFromDbRight, isExist := databaseTablesFromDbRight[k]

		// DB2 doesn't have table
		if !(isExist) {
			log.Print(fmt.Sprintf("Schema %s doesn't exist in compared db.", schemaFromDbLeft.SchemaName))
			if isSideLeft {
				comparerResult.LeftDatabaseExtraSchemas = append(comparerResult.LeftDatabaseExtraSchemas, schemaFromDbLeft)
			}

			if !(isSideLeft) {
				comparerResult.RightDatabaseExtraSchemas = append(comparerResult.RightDatabaseExtraSchemas, schemaFromDbLeft)
			}

			continue
		}

		compareTables(schemaFromDbLeft.Tables, schemaFromDbRight.Tables, comparerResult, isSideLeft)
	}
}

func compareTables(tablesLeft map[string]*model.Table, tablesRight map[string]*model.Table, comparerResult *ComparerResult, isSideLeft bool) {
	for tableName, tableStructFromLeft := range tablesLeft {
		tableStructFromDbRight, isExist := tablesRight[tableName]

		if !(isExist) {
			log.Print(fmt.Sprintf("Table %s doesn't exist in compared db.", tableStructFromLeft.TableName))

			if isSideLeft {
				comparerResult.TablesToInsertDBRight = append(comparerResult.TablesToInsertDBRight, tableStructFromLeft)
			}

			if !(isSideLeft) {
				comparerResult.TablesToInsertDBLeft = append(comparerResult.TablesToInsertDBLeft, tableStructFromLeft)
			}

			continue
		}

		compareColumns(tableStructFromLeft, tableStructFromDbRight, comparerResult, isSideLeft)
	}
}

func compareColumns(tableStructFromLeft *model.Table, tableStructFromDbRight *model.Table, comparerResult *ComparerResult, isSideLeft bool) {
	for columnName, columnFromLeft := range tableStructFromLeft.Columns {
		columnFromRight, isExist := tableStructFromDbRight.Columns[columnName]

		if !(isExist) {
			log.Print(fmt.Sprintf("Table %s doesn't exist in compared db.", tableStructFromLeft.TableName))

			if isSideLeft {
				comparerResult.ColumnToInsertDBLeft = append(comparerResult.ColumnToInsertDBLeft, columnFromLeft)
			}

			if !(isSideLeft) {
				comparerResult.ColumnToInsertDBRight = append(comparerResult.ColumnToInsertDBRight, columnFromLeft)
			}

			continue
		}

		fmt.Println(columnFromRight)
	}
}
