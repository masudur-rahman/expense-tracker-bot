package sqlite

import (
	"context"
	"database/sql"

	isql "github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/sqlite/lib"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	ctx       context.Context
	conn      *sql.Conn
	statement lib.Statement
}

func GetSqliteConnection() (*sql.Conn, error) {
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

func NewSqlite(ctx context.Context, conn *sql.Conn) Sqlite {
	return Sqlite{ctx: ctx, conn: conn}
}

func (pg Sqlite) Table(name string) isql.Database {
	pg.statement = pg.statement.Table(name)
	return pg
}

func (pg Sqlite) ID(id any) isql.Database {
	pg.statement = pg.statement.ID(id)
	return pg
}

func (pg Sqlite) In(col string, values ...any) isql.Database {
	pg.statement = pg.statement.In(col, values...)
	return pg
}

func (pg Sqlite) Where(cond string, args ...any) isql.Database {
	pg.statement = pg.statement.Where(cond, args...)
	return pg
}

func (pg Sqlite) Columns(cols ...string) isql.Database {
	pg.statement = pg.statement.Columns(cols...)
	return pg
}

func (pg Sqlite) AllCols() isql.Database {
	pg.statement = pg.statement.AllCols()
	return pg
}

func (pg Sqlite) MustCols(cols ...string) isql.Database {
	pg.statement = pg.statement.MustCols(cols...)
	return pg
}

func (pg Sqlite) ShowSQL(showSQL bool) isql.Database {
	pg.statement = pg.statement.ShowSQL(showSQL)
	return pg
}

func (pg Sqlite) FindOne(document any, filter ...any) (bool, error) {
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

func (pg Sqlite) FindMany(documents any, filter ...any) error {
	pg.statement = pg.statement.GenerateWhereClause(filter...)

	query := pg.statement.GenerateReadQuery()
	return pg.statement.ExecuteReadQuery(pg.ctx, pg.conn, query, documents)
}

func (pg Sqlite) InsertOne(document any) (id any, err error) {
	query := pg.statement.GenerateInsertQuery(document)
	return pg.statement.ExecuteInsertQuery(pg.ctx, pg.conn, query)
}

func (pg Sqlite) InsertMany(documents []any) ([]any, error) {
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

func (pg Sqlite) UpdateOne(document any) error {
	pg.statement = pg.statement.GenerateWhereClause()
	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := pg.statement.GenerateUpdateQuery(document)
	_, err := pg.statement.ExecuteWriteQuery(pg.ctx, pg.conn, query)
	return err
}

func (pg Sqlite) DeleteOne(filter ...any) error {
	pg.statement = pg.statement.GenerateWhereClause(filter...)
	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := pg.statement.GenerateDeleteQuery()
	_, err := pg.statement.ExecuteWriteQuery(pg.ctx, pg.conn, query)
	return err
}

func (pg Sqlite) Query(query string, args ...any) (*sql.Rows, error) {
	return pg.conn.QueryContext(pg.ctx, query, args...)
}

func (pg Sqlite) Exec(query string, args ...any) (sql.Result, error) {
	return pg.conn.ExecContext(pg.ctx, query, args...)
}

func (p Sqlite) Sync(tables ...any) error {
	ctx := context.Background()
	for _, table := range tables {
		if err := lib.SyncTable(ctx, p.conn, table); err != nil {
			return err
		}
	}

	return nil
}

func (pg Sqlite) cleanup() {
	pg.statement = lib.Statement{}
}
