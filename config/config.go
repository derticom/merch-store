package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

//nolint:tagliatelle // snake_case is allowed here.
type Config struct {
	LogLevel    string `yaml:"log_level"`
	PostgresURL string `yaml:"postgres_url"`
	Port        string `yaml:"port"`
}

func New() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
