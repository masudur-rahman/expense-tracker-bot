package sqlite

import (
	"context"
	"database/sql"
	"errors"

	isql "github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/sqlite/lib"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	ctx       context.Context
	conn      *sql.Conn
	tx        *sql.Tx
	statement lib.Statement
}

func NewSQLite(ctx context.Context, conn *sql.Conn) SQLite {
	return SQLite{ctx: ctx, conn: conn}
}

var _ isql.Database = SQLite{}

func (sq SQLite) BeginTx() (isql.Database, error) {
	if sq.tx != nil {
		return nil, errors.New("session already in progress")
	}
	tx, err := sq.conn.BeginTx(sq.ctx, nil)
	if err != nil {
		return nil, err
	}
	sq.tx = tx
	return sq, nil
}

func (sq SQLite) Commit() error {
	if sq.tx == nil {
		return errors.New("no transaction in progress")
	}
	err := sq.tx.Commit()
	sq.tx = nil
	return err
}

func (sq SQLite) Rollback() error {
	if sq.tx == nil {
		return errors.New("no transaction in progress")
	}
	err := sq.tx.Rollback()
	sq.tx = nil
	return err
}

func (sq SQLite) Table(name string) isql.Database {
	sq.statement = sq.statement.Table(name)
	return sq
}

func (sq SQLite) ID(id any) isql.Database {
	sq.statement = sq.statement.ID(id)
	return sq
}

func (sq SQLite) In(col string, values ...any) isql.Database {
	sq.statement = sq.statement.In(col, values...)
	return sq
}

func (sq SQLite) Where(cond string, args ...any) isql.Database {
	sq.statement = sq.statement.Where(cond, args...)
	return sq
}

func (sq SQLite) Columns(cols ...string) isql.Database {
	sq.statement = sq.statement.Columns(cols...)
	return sq
}

func (sq SQLite) AllCols() isql.Database {
	sq.statement = sq.statement.AllCols()
	return sq
}

func (sq SQLite) MustCols(cols ...string) isql.Database {
	sq.statement = sq.statement.MustCols(cols...)
	return sq
}

func (sq SQLite) ShowSQL(showSQL bool) isql.Database {
	sq.statement = sq.statement.ShowSQL(showSQL)
	return sq
}

func (sq SQLite) FindOne(document any, filter ...any) (bool, error) {
	sq.statement = sq.statement.GenerateWhereClause(filter...)

	if err := sq.statement.CheckWhereClauseNotEmpty(); err != nil {
		return false, err
	}

	query := sq.statement.GenerateReadQuery()
	err := sq.statement.ExecuteReadQuery(sq.ctx, sq.conn, sq.tx, query, document)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func (sq SQLite) FindMany(documents any, filter ...any) error {
	sq.statement = sq.statement.GenerateWhereClause(filter...)

	query := sq.statement.GenerateReadQuery()
	return sq.statement.ExecuteReadQuery(sq.ctx, sq.conn, sq.tx, query, documents)
}

func (sq SQLite) InsertOne(document any) (id any, err error) {
	query := sq.statement.GenerateInsertQuery(document)
	return sq.statement.ExecuteInsertQuery(sq.ctx, sq.conn, sq.tx, query)
}

func (sq SQLite) InsertMany(documents []any) ([]any, error) {
	var ids []any
	for _, doc := range documents {
		query := sq.statement.GenerateInsertQuery(doc)
		id, err := sq.statement.ExecuteInsertQuery(sq.ctx, sq.conn, sq.tx, query)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (sq SQLite) UpdateOne(document any) error {
	sq.statement = sq.statement.GenerateWhereClause()
	if err := sq.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := sq.statement.GenerateUpdateQuery(document)
	_, err := sq.statement.ExecuteWriteQuery(sq.ctx, sq.conn, sq.tx, query)
	return err
}

func (sq SQLite) DeleteOne(filter ...any) error {
	sq.statement = sq.statement.GenerateWhereClause(filter...)
	if err := sq.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := sq.statement.GenerateDeleteQuery()
	_, err := sq.statement.ExecuteWriteQuery(sq.ctx, sq.conn, sq.tx, query)
	return err
}

func (sq SQLite) Query(query string, args ...any) (*sql.Rows, error) {
	return sq.conn.QueryContext(sq.ctx, query, args...)
}

func (sq SQLite) Exec(query string, args ...any) (sql.Result, error) {
	return sq.conn.ExecContext(sq.ctx, query, args...)
}

func (sq SQLite) Sync(tables ...any) error {
	ctx := context.Background()
	for _, table := range tables {
		if err := lib.SyncTable(ctx, sq.conn, table); err != nil {
			return err
		}
	}

	return nil
}

func (sq SQLite) Close() error {
	return sq.conn.Close()
}

func (sq SQLite) cleanup() {
	sq.statement = lib.Statement{}
}
