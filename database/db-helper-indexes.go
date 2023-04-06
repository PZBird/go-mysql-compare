package db

import (
	"database/sql"
	"log"

	"github.com/PZBird/go-mysql-compare/model"
)

type indexResponse struct {
	IndexName  string
	SeqInIndex int8
	ColumnName string
	NonUnique  bool
	IndexType  string
	Comment    string
}

func readIndexes(conn *sql.DB, schema *model.DatabaseSchema,
	tableName string, table *model.Table, hostname string) {
	q := "SELECT"
	q += "	index_name,"
	q += "	seq_in_index,"
	q += "	column_name,"
	q += "	non_unique,"
	q += "	index_type,"
	q += "	comment"
	q += " FROM"
	q += "		INFORMATION_SCHEMA.STATISTICS"
	q += " WHERE 1=1"
	q += "		AND table_schema = ?"
	q += "		AND table_name = ?"
	q += " ORDER BY seq_in_index;"

	table.Indexes = make(map[string]*model.Index)

	rows, err := conn.Query(q, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		index := &model.Index{}
		index.TableName = tableName
		index.DatabaseName = schema.SchemaName

		var response indexResponse

		err := rows.Scan(
			&response.IndexName,
			&response.SeqInIndex,
			&response.ColumnName,
			&response.NonUnique,
			&response.IndexType,
			&response.Comment,
		)

		if err != nil {
			log.Fatal(err)
		}

		appendIndex(index, response, schema, tableName)
	}
}

func appendIndex(index *model.Index, response indexResponse, schema *model.DatabaseSchema, tableName string) {
	indexes := schema.Tables[tableName].Indexes

	if indexes[response.IndexName] != nil {
		indexes[response.IndexName].Columns = append(indexes[response.IndexName].Columns, &model.IndexColumn{SeqInIndex: response.SeqInIndex, ColumnName: response.ColumnName})
		return
	}
	index.IndexName = response.IndexName
	index.Columns = append(index.Columns, &model.IndexColumn{SeqInIndex: response.SeqInIndex, ColumnName: response.ColumnName})
	index.Comment = response.Comment
	index.NonUnique = response.NonUnique

	indexes[response.IndexName] = index
}
