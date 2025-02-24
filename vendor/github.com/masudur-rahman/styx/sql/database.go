package sql

import (
	"database/sql"
)

type Engine interface {
	BeginTx() (Engine, error)
	Commit() error
	Rollback() error

	Table(name string) Engine

	ID(id any) Engine
	In(col string, values ...any) Engine
	Where(cond string, args ...any) Engine
	Columns(cols ...string) Engine
	AllCols() Engine
	MustCols(cols ...string) Engine
	ShowSQL(showSQL bool) Engine

	FindOne(document any, filter ...any) (bool, error)
	FindMany(documents any, filter ...any) error

	InsertOne(document any) (id any, err error)
	// TODO: might need to convert []any to just any
	InsertMany(documents []any) ([]any, error)

	UpdateOne(document any) error

	DeleteOne(filter ...any) error

	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)

	Sync(tables ...any) error

	Close() error
}
