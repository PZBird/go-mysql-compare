package model

import "database/sql"

type DatabaseSchema struct {
	DatabaseName string
	SchemaName   string
	Tables       map[string]*Table
}

type Table struct {
	DatabaseName string
	TableName    string
	Columns      map[string]*Column
	PrimaryKeys  []*Column
	OtherColumns []*Column
}

type View struct {
	DatabaseName string
	ViewName     string
	Columns      []*Column
}

type Column struct {
	DatabaseName           string
	TableName              string
	ColumnName             string
	IsPrimaryKey           bool
	IsNullable             bool
	IsAutoIncrement        bool
	IsUnique               bool
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	DefaultValue           string
	EnumValues             []string
}
