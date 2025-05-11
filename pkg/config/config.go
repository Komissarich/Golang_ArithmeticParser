package config

import (
<<<<<<< HEAD
	"calc/server/models"
	"time"
=======
	"calc/server/agent"
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
<<<<<<< HEAD
	Agent          models.AgentConfig `yaml:"AGENT"`
	Server_Port    string             `yaml:"SERVER_PORT"`
	Jwt_secret_key string             `yaml:"JWT_SECRET"`
	Jwt_expiration time.Duration      `yaml:"JWT_EXPIRATION"`
=======
	Agent       agent.Config `yaml:"AGENT"`
	Server_Port string       `yaml:"SERVER_PORT"`
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
