package sqlite

import (
	"context"
	"database/sql"

	isql "github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/sqlite/lib"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	ctx       context.Context
	conn      *sql.Conn
	statement lib.Statement
}

func GetSQLiteConnection() (*sql.Conn, error) {
	db, err := sql.Open("sqlite3", "expense-tracker.db")
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	if err = conn.PingContext(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}

func NewSQLite(ctx context.Context, conn *sql.Conn) SQLite {
	return SQLite{ctx: ctx, conn: conn}
}

func (pg SQLite) Table(name string) isql.Database {
	pg.statement = pg.statement.Table(name)
	return pg
}

func (pg SQLite) ID(id any) isql.Database {
	pg.statement = pg.statement.ID(id)
	return pg
}

func (pg SQLite) In(col string, values ...any) isql.Database {
	pg.statement = pg.statement.In(col, values...)
	return pg
}

func (pg SQLite) Where(cond string, args ...any) isql.Database {
	pg.statement = pg.statement.Where(cond, args...)
	return pg
}

func (pg SQLite) Columns(cols ...string) isql.Database {
	pg.statement = pg.statement.Columns(cols...)
	return pg
}

func (pg SQLite) AllCols() isql.Database {
	pg.statement = pg.statement.AllCols()
	return pg
}

func (pg SQLite) MustCols(cols ...string) isql.Database {
	pg.statement = pg.statement.MustCols(cols...)
	return pg
}

func (pg SQLite) ShowSQL(showSQL bool) isql.Database {
	pg.statement = pg.statement.ShowSQL(showSQL)
	return pg
}

func (pg SQLite) FindOne(document any, filter ...any) (bool, error) {
	pg.statement = pg.statement.GenerateWhereClause(filter...)

	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return false, err
	}

	query := pg.statement.GenerateReadQuery()
	err := pg.statement.ExecuteReadQuery(pg.ctx, pg.conn, query, document)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func (pg SQLite) FindMany(documents any, filter ...any) error {
	pg.statement = pg.statement.GenerateWhereClause(filter...)

	query := pg.statement.GenerateReadQuery()
	return pg.statement.ExecuteReadQuery(pg.ctx, pg.conn, query, documents)
}

func (pg SQLite) InsertOne(document any) (id any, err error) {
	query := pg.statement.GenerateInsertQuery(document)
	return pg.statement.ExecuteInsertQuery(pg.ctx, pg.conn, query)
}

func (pg SQLite) InsertMany(documents []any) ([]any, error) {
	var ids []any
	for _, doc := range documents {
		query := pg.statement.GenerateInsertQuery(doc)
		id, err := pg.statement.ExecuteInsertQuery(pg.ctx, pg.conn, query)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (pg SQLite) UpdateOne(document any) error {
	pg.statement = pg.statement.GenerateWhereClause()
	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := pg.statement.GenerateUpdateQuery(document)
	_, err := pg.statement.ExecuteWriteQuery(pg.ctx, pg.conn, query)
	return err
}

func (pg SQLite) DeleteOne(filter ...any) error {
	pg.statement = pg.statement.GenerateWhereClause(filter...)
	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := pg.statement.GenerateDeleteQuery()
	_, err := pg.statement.ExecuteWriteQuery(pg.ctx, pg.conn, query)
	return err
}

func (pg SQLite) Query(query string, args ...any) (*sql.Rows, error) {
	return pg.conn.QueryContext(pg.ctx, query, args...)
}

func (pg SQLite) Exec(query string, args ...any) (sql.Result, error) {
	return pg.conn.ExecContext(pg.ctx, query, args...)
}

func (p SQLite) Sync(tables ...any) error {
	ctx := context.Background()
	for _, table := range tables {
		if err := lib.SyncTable(ctx, p.conn, table); err != nil {
			return err
		}
	}

	return nil
}

func (pg SQLite) cleanup() {
	pg.statement = lib.Statement{}
}
