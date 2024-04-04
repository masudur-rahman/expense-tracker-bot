package sql

import (
	"database/sql"
)

type Database interface {
	BeginTx() (Database, error)
	Commit() error
	Rollback() error

	Table(name string) Database

	ID(id any) Database
	In(col string, values ...any) Database
	Where(cond string, args ...any) Database
	Columns(cols ...string) Database
	AllCols() Database
	MustCols(cols ...string) Database
	ShowSQL(showSQL bool) Database

	FindOne(document any, filter ...any) (bool, error)
	FindMany(documents any, filter ...any) error

	InsertOne(document any) (id any, err error)
	InsertMany(documents []any) ([]any, error)

	UpdateOne(document any) error

	DeleteOne(filter ...any) error

	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)

	Sync(tables ...any) error

	Close() error
}
