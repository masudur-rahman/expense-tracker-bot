package postgres

import (
	"context"
	"database/sql"

	isql "github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/postgres/lib"
)

type Postgres struct {
	ctx       context.Context
	conn      *sql.Conn
	statement lib.Statement
}

func NewPostgres(ctx context.Context, conn *sql.Conn) Postgres {
	return Postgres{ctx: ctx, conn: conn}
}

func (pg Postgres) Table(name string) isql.Database {
	pg.statement = pg.statement.Table(name)
	return pg
}

func (pg Postgres) ID(id any) isql.Database {
	pg.statement = pg.statement.ID(id)
	return pg
}

func (pg Postgres) In(col string, values ...any) isql.Database {
	pg.statement = pg.statement.In(col, values...)
	return pg
}

func (pg Postgres) Where(cond string, args ...any) isql.Database {
	pg.statement = pg.statement.Where(cond, args...)
	return pg
}

func (pg Postgres) Columns(cols ...string) isql.Database {
	pg.statement = pg.statement.Columns(cols...)
	return pg
}

func (pg Postgres) AllCols() isql.Database {
	pg.statement = pg.statement.AllCols()
	return pg
}

func (pg Postgres) MustCols(cols ...string) isql.Database {
	pg.statement = pg.statement.MustCols(cols...)
	return pg
}

func (pg Postgres) ShowSQL(showSQL bool) isql.Database {
	pg.statement = pg.statement.ShowSQL(showSQL)
	return pg
}

func (pg Postgres) FindOne(document any, filter ...any) (bool, error) {
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

func (pg Postgres) FindMany(documents any, filter ...any) error {
	pg.statement = pg.statement.GenerateWhereClause(filter...)

	query := pg.statement.GenerateReadQuery()
	return pg.statement.ExecuteReadQuery(pg.ctx, pg.conn, query, documents)
}

func (pg Postgres) InsertOne(document any) (id any, err error) {
	query := pg.statement.GenerateInsertQuery(document)
	return pg.statement.ExecuteInsertQuery(pg.ctx, pg.conn, query)
}

func (pg Postgres) InsertMany(documents []any) ([]any, error) {
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

func (pg Postgres) UpdateOne(document any) error {
	pg.statement = pg.statement.GenerateWhereClause()
	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := pg.statement.GenerateUpdateQuery(document)
	_, err := pg.statement.ExecuteWriteQuery(pg.ctx, pg.conn, query)
	return err
}

func (pg Postgres) DeleteOne(filter ...any) error {
	pg.statement = pg.statement.GenerateWhereClause(filter...)
	if err := pg.statement.CheckWhereClauseNotEmpty(); err != nil {
		return err
	}

	query := pg.statement.GenerateDeleteQuery()
	_, err := pg.statement.ExecuteWriteQuery(pg.ctx, pg.conn, query)
	return err
}

func (pg Postgres) Query(query string, args ...any) (*sql.Rows, error) {
	return pg.conn.QueryContext(pg.ctx, query, args...)
}

func (pg Postgres) Exec(query string, args ...any) (sql.Result, error) {
	return pg.conn.ExecContext(pg.ctx, query, args...)
}

func (pg Postgres) Sync(tables ...any) error {
	ctx := context.Background()
	for _, table := range tables {
		if err := lib.SyncTable(ctx, pg.conn, table); err != nil {
			return err
		}
	}

	return nil
}

func (p Postgres) Close() error {
	return p.conn.Close()
}

func (pg Postgres) cleanup() {
	pg.statement = lib.Statement{}
}
