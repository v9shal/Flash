package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	RedisURL    string `env:"REDISURL,required"`
	JwtSecret   string `env:"JWTSECRET,required"`
}

func Load() (*Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
