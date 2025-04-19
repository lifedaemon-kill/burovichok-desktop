package config

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"
)

// Единое место для хранения ключей конфигурации.

const (
	PathConfig = "config/config.yaml"
)

type Config struct {
	ENV    string     `yaml:"env" env-required:"true"`
	DB     DBConf     `yaml:"db" env-required:"true"`
	Logger LoggerConf `yaml:"logger" env-required:"true"`
	UI     UI         `yaml:"ui" env-required:"true"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "no such file %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}
	return &cfg, nil
}

type DBConf struct {
	DSN             string `yaml:"dsn"`
	MigrationsPath  string `yaml:"confmigration_path"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_Idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	MaxRetries      int    `yaml:"max_retries"`
}

type LoggerConf struct {
	Env string `yaml:"env" env-required:"true"`
}

type UI struct {
	Name     string `yaml:"name" env-required:"true"`
	Width    int    `yaml:"width" env-required:"true"`
	Height   int    `yaml:"height" env-required:"true"`
	IconPath string `yaml:"icon_path" env-required:"true"`
}
