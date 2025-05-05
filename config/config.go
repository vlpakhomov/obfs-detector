package config

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
)

type Config struct {
	Postgres Postgres `validate:"required"`
}

type Postgres struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
}

func Load(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open config file fail: %v\n", err)
	}

	cfg := Config{}

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("decode config fail: %v\n", err)
	}

	if err = validator.New().Struct(&cfg); err != nil {
		log.Fatalf("validate config fail: %v\n", err)
	}

	return cfg
}
