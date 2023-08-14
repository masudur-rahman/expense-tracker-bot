/*
Copyright © 2023 Masudur Rahman <masudjuly02@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	isql "github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/expense-tracker-bot/api"
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/masudur-rahman/database/sql/postgres"
	"github.com/masudur-rahman/database/sql/postgres/lib"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		go startHealthz()

		if err := getServicesForPostgres(cmd.Context()); err != nil {
			log.Fatalln(err)
		}

		bot, err := api.TeleBotRoutes()
		if err != nil {
			log.Fatalln(err)
		}

		go pingHealthzApiPeriodically()
		log.Println("Expense Tracker Bot started")
		bot.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func startHealthz() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Running"))
	})

	logr.DefaultLogger.Infow("Health checker started at :8080/healthz")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}

//func getServicesForSupabase(ctx context.Context) *all.Services {
//	supClient := supabase.InitializeSupabase(ctx)
//
//	var db isql.Database
//	db = supabase.NewSupabase(ctx, supClient)
//	logger := logr.DefaultLogger
//	return all.InitiateSQLServices(db, logger)
//}

func getServicesForPostgres(ctx context.Context) error {
	cfg := parsePostgresConfig()

	err := initiateSQLServices(ctx, cfg)
	if err != nil {
		return err
	}

	if err = all.GetServices().Txn.UpdateTxnCategories(); err != nil {
		return err
	}
	return nil
}

func initiateSQLServices(ctx context.Context, cfg lib.PostgresConfig) error {
	conn, err := lib.GetPostgresConnection(cfg)
	if err != nil {
		return err
	}

	db := postgres.NewPostgres(ctx, conn).ShowSQL(true)
	syncTables(db)

	logger := logr.DefaultLogger
	all.InitiateSQLServices(db, logger)

	go pingPostgresDatabasePeriodically(ctx, cfg, conn, logger)

	return nil
}

func pingPostgresDatabasePeriodically(ctx context.Context, cfg lib.PostgresConfig, conn *sql.Conn, logger logr.Logger) {
	t5 := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-t5.C:
			if err := conn.PingContext(ctx); err != nil {
				logger.Errorw("Database connection closed", "error", err.Error())
				conn, err = lib.GetPostgresConnection(cfg)
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

func pingHealthzApiPeriodically() {
	logger := logr.DefaultLogger
	baseURL, ok := os.LookupEnv("BASE_URL")
	if !ok {
		return
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalln(err)
	}
	u.Path = path.Join(u.Path, "healthz")
	healthPath := u.String()
	logger.Infow("Health url provided", "url", healthPath)

	t20 := time.NewTicker(20 * time.Minute)
	for {
		select {
		case <-t20.C:
			resp, err := http.Get(healthPath)
			if err != nil {
				logger.Errorw("healthz api failed", "error", err.Error())
			} else {
				data, err := io.ReadAll(resp.Body)
				var errMsg string
				if err != nil {
					errMsg = err.Error()
				}
				logger.Infow("healthz api", "status", resp.StatusCode, "msg", string(data), "error", errMsg)
			}
		}
	}
}

func parsePostgresConfig() lib.PostgresConfig {
	cfg := lib.PostgresConfig{
		Name:     "expense",
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		SSLMode:  "disable",
	}

	user, ok := os.LookupEnv("POSTGRES_USER")
	if ok {
		cfg.User = user
	}
	pass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if ok {
		cfg.Password = pass
	}
	name, ok := os.LookupEnv("POSTGRES_DB")
	if ok {
		cfg.Name = name
	}
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if ok {
		cfg.Host = host
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if ok {
		cfg.Port = port
	}
	ssl, ok := os.LookupEnv("POSTGRES_SSL_MODE")
	if ok {
		cfg.SSLMode = ssl
	}
	return cfg
}

func syncTables(db isql.Database) {
	err := db.Sync(
		models.User{},
		models.Account{},
		models.Transaction{},
		models.TxnCategory{},
		models.TxnSubcategory{},
		models.Event{},
	)
	if err != nil {
		log.Fatalln(err)
	}
}
