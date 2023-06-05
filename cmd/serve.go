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
	"log"

	"github.com/masudur-rahman/expense-tracker-bot/api"
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/supabase"

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
		svc := getServicesForSupabase(cmd.Context())
		bot, err := api.TeleBotRoutes(svc)
		if err != nil {
			panic(err)
		}

		log.Println("Expense Tracker Bot started")
		bot.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func getServicesForSupabase(ctx context.Context) *all.Services {
	supClient := supabase.InitializeSupabase(ctx)

	var db sql.Database
	db = supabase.NewSupabase(ctx, supClient)
	logger := logr.DefaultLogger
	return all.GetSQLServices(db, logger)
}
