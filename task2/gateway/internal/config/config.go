package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string          `yaml:"env" env-required:"true"`
	LogLevel string          `yaml:"log-level" env-default:"info"`
	HTTP     HTTPConfig      `yaml:"http"`
	GRPC     CollectorConfig `yaml:"grpc"`
}

type HTTPConfig struct {
	Port           string `yaml:"port" env-default:":8080"`
	TimeoutSeconds int    `yaml:"timeout_seconds" env-default:"30"`
	Timeout        time.Duration
}

type CollectorConfig struct {
	Address        string `yaml:"address" env-required:"true"`
	TimeoutSeconds int    `yaml:"timeout_seconds" env-default:"5"`
	Timeout        time.Duration
}

func Load() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg.GRPC.Timeout = time.Duration(cfg.GRPC.TimeoutSeconds) * time.Second
	cfg.HTTP.Timeout = time.Duration(cfg.HTTP.TimeoutSeconds) * time.Second

	return &cfg, nil
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
