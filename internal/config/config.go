package config

import "github.com/caarlos0/env/v7"

type Config struct {
	ServerPort    int    `env:"SERVER_PORT" envDefault:"8080"`
	Neo4jURL      string `env:"NEO4J_URL,notEmpty"`
	Neo4jUsername string `env:"NEO4J_USERNAME,notEmpty"`
	Neo4jPassword string `env:"NEO4J_PASSWORD,notEmpty"`
}

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
