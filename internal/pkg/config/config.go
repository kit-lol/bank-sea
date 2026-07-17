package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config describes the main configuration structure
type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HttpServer HTTPServer `yaml:"http_server"`
	DB         Database   `yaml:"db"`
}

// HTTPServer describes the web server settings
type HTTPServer struct {
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Database describes the database connection parameters
type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

// MustLoad reads the YAML file and fills in the Config structure
func MustLoad() *Config {
	configPath := "configs/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("The configuration file was not found on the path: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error when reading the configuration file: %s", err)
	}

	return &cfg
}
