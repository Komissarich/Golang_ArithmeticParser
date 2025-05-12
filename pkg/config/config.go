package config

import (
	"calc/server/models"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Agent          models.AgentConfig `yaml:"AGENT"`
	Server_Port    string             `yaml:"SERVER_PORT"`
	Jwt_secret_key string             `yaml:"JWT_SECRET"`
	Jwt_expiration time.Duration      `yaml:"JWT_EXPIRATION"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
