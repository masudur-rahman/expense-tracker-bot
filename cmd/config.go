package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/masudur-rahman/expense-tracker-bot/configs"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"

	"gopkg.in/yaml.v3"
)

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		cfgFile = filepath.Join(pkg.ProjectDirectory, "configs", ".expense-tracker.yaml")
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("Reading config file %v, %v", cfgFile, err)
	}

	if err = yaml.Unmarshal(data, &configs.TrackerConfig); err != nil {
		log.Fatalf("Unmarshaling PurrfectConfig, %v", err)
	}
}
