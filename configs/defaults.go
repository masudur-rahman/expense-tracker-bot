package configs

import (
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
)

var (
	once    sync.Once
	printer pkg.Printer
)

func GetDefaultPrinter() pkg.Printer {
	once.Do(func() {
		opts := pkg.Options{Style: table.StyleColoredBright, EnableStdout: true}
		printer = pkg.NewPrinter(opts)
	})
	return printer
}
