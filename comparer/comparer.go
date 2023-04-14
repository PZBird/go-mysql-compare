package comparer

import (
	"fmt"
	"log"
	"reflect"

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
	ColumnToRemoveDBLeft      []*model.Column
	ColumnToRemoveDBRight     []*model.Column
	ColumnToModifyDBLeft      []*model.Column
	ColumnToModifyDBRight     []*model.Column
	EnumToModifyDBLeft        []*model.Column
	EnumToModifyDBRight       []*model.Column
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
	log.Printf("Schema %s doesn't exist in compared db.", schema.SchemaName)

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
	log.Printf("Index %s doesn't exist in compared %s.", index.IndexName, tableStruct.TableName)

	if isSideLeft {
		comparerResult.IndexToInsertDBLeft = append(comparerResult.IndexToInsertDBLeft, index)
		return
	}

	comparerResult.IndexToInsertDBRight = append(comparerResult.IndexToInsertDBRight, index)
}

func appendTableForInsert(tableStruct *model.Table, isSideLeft bool, comparerResult *ComparerResult) {
	log.Printf("Table %s doesn't exist in compared db.", tableStruct.TableName)

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

		deepColumnCompare(columnFromLeft, columnFromRight, isSideLeft, comparerResult)

		fmt.Println(columnFromRight)
	}
}

func deepColumnCompare(columnFromLeft *model.Column, columnFromRight *model.Column, isSideLeft bool, comparerResult *ComparerResult) {
	if len(columnFromLeft.EnumValues) > 0 {
		compareEnums(columnFromLeft, columnFromRight, isSideLeft, comparerResult)
	}

	back := columnFromLeft.DatabaseName
	columnFromLeft.DatabaseName = columnFromRight.DatabaseName

	if !reflect.DeepEqual(columnFromLeft, columnFromRight) {
		comparerResult.ColumnToModifyDBLeft = append(comparerResult.ColumnToModifyDBLeft, columnFromRight)
		comparerResult.ColumnToModifyDBRight = append(comparerResult.ColumnToModifyDBRight, columnFromLeft)
	}

	columnFromLeft.DatabaseName = back
}

func compareEnums(columnFromLeft, columnFromRight *model.Column, isSideLeft bool, comparerResult *ComparerResult) {
	if isSideLeft {
		return
	}

	if len(columnFromLeft.EnumValues) != len(columnFromRight.EnumValues) {
		comparerResult.EnumToModifyDBLeft = append(comparerResult.EnumToModifyDBLeft, columnFromRight)
		comparerResult.EnumToModifyDBRight = append(comparerResult.EnumToModifyDBRight, columnFromLeft)

		return
	}

	for _, value := range columnFromLeft.EnumValues {
		if !contains(columnFromRight.EnumValues, value) {
			comparerResult.EnumToModifyDBLeft = append(comparerResult.EnumToModifyDBLeft, columnFromRight)
			comparerResult.EnumToModifyDBRight = append(comparerResult.EnumToModifyDBRight, columnFromLeft)

			return
		}
	}
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}

	return false
}

func appendColumnsForDb(tableStruct *model.Table, isSideLeft bool, comparerResult *ComparerResult, column *model.Column) {
	log.Printf("Column %s doesn't exist in compared %s.", column.ColumnName, tableStruct.TableName)

	if isSideLeft {
		comparerResult.ColumnToInsertDBLeft = append(comparerResult.ColumnToInsertDBLeft, column)
		comparerResult.ColumnToRemoveDBRight = append(comparerResult.ColumnToRemoveDBRight, column)
		return
	}

	comparerResult.ColumnToInsertDBRight = append(comparerResult.ColumnToInsertDBRight, column)
	comparerResult.ColumnToRemoveDBLeft = append(comparerResult.ColumnToRemoveDBLeft, column)
}
