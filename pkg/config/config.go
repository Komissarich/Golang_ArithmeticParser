package config

import (
	"calc/server/agent"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Agent       agent.Config `yaml:"AGENT"`
	Server_Port string       `yaml:"SERVER_PORT"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
