package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"klavio-template/internal/app/config"
	"klavio-template/internal/app/scanner"
	"log"
	"os"
)

// configPath flag determines absolute path to config file.
var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "./main.yml", "Path to config file.")
	flag.Parse()
}

func main() {
	cfg := getConfig()

	s := scanner.NewScanner(cfg)
	if err := s.Start(); err != nil {
		log.Fatalf("failed to start the scan: %v", err)
	}
}

// getConfig returns new App.
func getConfig() config.App {
	b, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}

	cfg := config.New()
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatalf("failed to unmarshall config file: %v", err)
	}

	return cfg
}
