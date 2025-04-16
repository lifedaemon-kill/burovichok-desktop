package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	ENV    string     `yaml:"env" env-required:"true"`
	DB     DBConf     `yaml:"db" env-required:"true"`
	Logger LoggerConf `yaml:"logger" env-required:"true"`
}

func Load(configPath string) *Config {
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("no such file ", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}

type DBConf struct {
	dsn string `yaml:"dsn" env-required:"true"`
}
type LoggerConf struct {
	Env string `yaml:"env" env-required:"true"`
}
