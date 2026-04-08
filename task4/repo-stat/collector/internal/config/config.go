package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string       `yaml:"log-level" env-default:"INFO"`
	GRPC     GRPCServer   `yaml:"grpc"`
	Github   GithubConfig `yaml:"github"`
}

type GRPCServer struct {
	Port           string `yaml:"port"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	Timeout        time.Duration
}

type GithubConfig struct {
	BaseURL        string `yaml:"baseurl"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	Timeout        time.Duration
	UserAgent      string `yaml:"user_agent" env-default:"Collector-Service/1.0"`
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
	cfg.Github.Timeout = time.Duration(cfg.Github.TimeoutSeconds) * time.Second

	return &cfg, nil
}

func fetchConfigPath() string {
	var res string

	f := flag.Lookup("config")
	if f != nil {
		res = f.Value.String()
	} else {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
