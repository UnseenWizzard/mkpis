package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type configuration struct {
	GitHubToken string `env:"GITHUB_TOKEN"`
}

func loadConfig() *configuration {
	godotenv.Load() //load .env

	cfg := &configuration{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	return cfg
}

var Env = loadConfig()
