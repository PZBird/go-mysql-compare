package comparer

import (
	"fmt"
	"log"

	"github.com/PZBird/go-mysql-compare/configuration"
	"github.com/PZBird/go-mysql-compare/model"
)

type ComparerResult struct {
	LeftDatabaseExtraSchemas  []*model.DatabaseSchema
	RightDatabaseExtraSchemas []*model.DatabaseSchema
	TablesToInsertDBLeft      []*model.Table
	TablesToInsertDBRight     []*model.Table
	ColumnToInsertDBLeft      []*model.Column
	ColumnToInsertDBRight     []*model.Column
	IndexToInsertDBLeft       []*model.Index
	IndexToInsertDBRight      []*model.Index
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
			appendExtraSchema(isSideLeft, comparerResult, schemaFromDbLeft)

			continue
		}

		compareTables(schemaFromDbLeft.Tables, schemaFromDbRight.Tables, comparerResult, isSideLeft)
	}
}

func appendExtraSchema(isSideLeft bool, comparerResult *ComparerResult, schema *model.DatabaseSchema) {
	log.Print(fmt.Sprintf("Schema %s doesn't exist in compared db.", schema.SchemaName))

	if isSideLeft {
		comparerResult.LeftDatabaseExtraSchemas = append(comparerResult.LeftDatabaseExtraSchemas, schema)
		return
	}

	comparerResult.RightDatabaseExtraSchemas = append(comparerResult.RightDatabaseExtraSchemas, schema)
}

func compareTables(tablesLeft map[string]*model.Table, tablesRight map[string]*model.Table, comparerResult *ComparerResult, isSideLeft bool) {
	for tableName, tableStructFromLeft := range tablesLeft {
		tableStructFromDbRight, isExist := tablesRight[tableName]

		if !(isExist) {
			appendTableForInsert(tableStructFromLeft, isSideLeft, comparerResult)

			continue
		}

		compareColumns(tableStructFromLeft, tableStructFromDbRight, comparerResult, isSideLeft)
		compareIndexes(tableStructFromLeft, tableStructFromDbRight, comparerResult, isSideLeft)
	}
}

func compareIndexes(tableStruct *model.Table, tableStructForCompare *model.Table, comparerResult *ComparerResult, isSideLeft bool) {
	for indexName, indexOrigin := range tableStruct.Indexes {
		indexFromRight, isExist := tableStructForCompare.Indexes[indexName]

		if !(isExist) {
			appendIndexForTable(tableStruct, isSideLeft, comparerResult, indexOrigin)

			continue
		}

		fmt.Println(indexFromRight)
	}
}

func appendIndexForTable(tableStruct *model.Table, isSideLeft bool, comparerResult *ComparerResult, index *model.Index) {
	log.Print(fmt.Sprintf("Index %s doesn't exist in compared %s.", index.IndexName, tableStruct.TableName))

	if isSideLeft {
		comparerResult.IndexToInsertDBLeft = append(comparerResult.IndexToInsertDBLeft, index)
		return
	}

	comparerResult.IndexToInsertDBRight = append(comparerResult.IndexToInsertDBRight, index)
}

func appendTableForInsert(tableStruct *model.Table, isSideLeft bool, comparerResult *ComparerResult) {
	log.Print(fmt.Sprintf("Table %s doesn't exist in compared db.", tableStruct.TableName))

	if isSideLeft {
		comparerResult.TablesToInsertDBRight = append(comparerResult.TablesToInsertDBRight, tableStruct)
		return
	}

	comparerResult.TablesToInsertDBLeft = append(comparerResult.TablesToInsertDBLeft, tableStruct)
}

func compareColumns(tableStruct *model.Table, tableStructForCompare *model.Table, comparerResult *ComparerResult, isSideLeft bool) {
	for columnName, columnFromLeft := range tableStruct.Columns {
		columnFromRight, isExist := tableStructForCompare.Columns[columnName]

		if !(isExist) {
			appendColumnsForDb(tableStruct, isSideLeft, comparerResult, columnFromLeft)

			continue
		}

		fmt.Println(columnFromRight)
	}
}

func appendColumnsForDb(tableStruct *model.Table, isSideLeft bool, comparerResult *ComparerResult, column *model.Column) {
	log.Print(fmt.Sprintf("Column %s doesn't exist in compared %s.", column.ColumnName, tableStruct.TableName))

	if isSideLeft {
		comparerResult.ColumnToInsertDBLeft = append(comparerResult.ColumnToInsertDBLeft, column)
		return
	}

	comparerResult.ColumnToInsertDBRight = append(comparerResult.ColumnToInsertDBRight, column)
}
