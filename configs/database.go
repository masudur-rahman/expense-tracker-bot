package configs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/modules/google"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	isql "github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/postgres"
	sqlib "github.com/masudur-rahman/database/sql/postgres/lib"
	"github.com/masudur-rahman/database/sql/sqlite"
	"github.com/masudur-rahman/database/sql/sqlite/lib"
)

func InitiateDatabaseConnection(ctx context.Context) error {
	cfg := TrackerConfig.Database
	switch cfg.Type {
	case DatabasePostgres:
		db, err := getPostgresDatabase(ctx)
		if err != nil {
			return err
		}
		return initializeSQLServices(db)
	case DatabaseSqlite, "":
		if cfg.Sqlite.SyncToDrive {
			if !cfg.Sqlite.DisableSyncFromDrive {
				if err := google.SyncDatabaseFromDrive(); err != nil {
					return err
				}
				logr.DefaultLogger.Infof("Sqlite database synced from google drive")
			}
			go google.SyncDatabaseToDrivePeriodically(TrackerConfig.Database.Sqlite.SyncInterval)
		}

		db, err := getSqliteDatabase(ctx)
		if err != nil {
			return err
		}
		return initializeSQLServices(db)
	default:
		return fmt.Errorf("unknown database type")
	}
}

func getSqliteDatabase(ctx context.Context) (isql.Database, error) {
	conn, err := lib.GetSQLiteConnection("expense-tracker.db")
	if err != nil {
		return nil, err
	}

	return sqlite.NewSqlite(ctx, conn), nil
}

func initializeSQLServices(db isql.Database) error {
	if err := syncTables(db); err != nil {
		return err
	}
	all.InitiateSQLServices(db, logr.DefaultLogger)

	return all.GetServices().Txn.UpdateTxnCategories()
}

//func getServicesForSupabase(ctx context.Context) *all.Services {
//	supClient := supabase.InitializeSupabase(ctx)
//
//	var db isql.Database
//	db = supabase.NewSupabase(ctx, supClient)
//	logger := logr.DefaultLogger
//	return all.InitiateSQLServices(db, logger)
//}

func getPostgresDatabase(ctx context.Context) (isql.Database, error) {
	parsePostgresConfig()
	conn, err := sqlib.GetPostgresConnection(TrackerConfig.Database.Postgres)
	if err != nil {
		return nil, err
	}
	go pingPostgresDatabasePeriodically(ctx, TrackerConfig.Database.Postgres, conn, logr.DefaultLogger)

	return postgres.NewPostgres(ctx, conn).ShowSQL(true), nil
}

func pingPostgresDatabasePeriodically(ctx context.Context, cfg sqlib.PostgresConfig, conn *sql.Conn, logger logr.Logger) {
	t5 := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-t5.C:
			if err := conn.PingContext(ctx); err != nil {
				logger.Errorw("Database connection closed", "error", err.Error())
				conn, err = sqlib.GetPostgresConnection(cfg)
				if err != nil {
					logger.Errorw("couldn't create database connection", "error", err.Error())
				}

				db := postgres.NewPostgres(ctx, conn).ShowSQL(true)
				all.InitiateSQLServices(db, logger)
				logger.Infow("New connection established")
			}
		}
	}
}

func parsePostgresConfig() {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if ok {
		TrackerConfig.Database.Postgres.User = user
	}
	pass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if ok {
		TrackerConfig.Database.Postgres.Password = pass
	}
	name, ok := os.LookupEnv("POSTGRES_DB")
	if ok {
		TrackerConfig.Database.Postgres.Name = name
	}
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if ok {
		TrackerConfig.Database.Postgres.Host = host
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if ok {
		TrackerConfig.Database.Postgres.Port = port
	}
	ssl, ok := os.LookupEnv("POSTGRES_SSL_MODE")
	if ok {
		TrackerConfig.Database.Postgres.SSLMode = ssl
	}
}

func syncTables(db isql.Database) error {
	return db.Sync(
		models.User{},
		models.Account{},
		models.Transaction{},
		models.TxnCategory{},
		models.TxnSubcategory{},
		models.Event{},
	)
}