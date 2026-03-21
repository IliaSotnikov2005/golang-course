package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string     `yaml:"env" env-required:"true"`
	LogLevel string     `yaml:"log-level" env-default:"info"`
	GRPC     GRPCConfig `yaml:"grpc"`
	HTTP     HTTPConfig `yaml:"http"`
}

type GRPCConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPConfig struct {
	BaseURL   string        `yaml:"baseurl"`
	Timeout   time.Duration `yaml:"timeout"`
	UserAgent string        `yaml:"user_agent" env-default:"Collector-Service/1.0"`
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
