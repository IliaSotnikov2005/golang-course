package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel    string      `yaml:"log-level" env-default:"INFO"`
	GRPCServer  GRPCServer  `yaml:"grpc_server"`
	KafkaConfig KafkaConfig `yaml:"kafka"`
	DatabaseDSN string      `yaml:"database_dsn"`
}

type GRPCServer struct {
	Port           string `yaml:"port"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	Timeout        time.Duration
}

type KafkaConfig struct {
	Brokers       []string `yaml:"brokers" env-required:"true"`
	RequestTopic  string   `yaml:"requests_topic" env-default:"repository-request"`
	ResponseTopic string   `yaml:"response_topic" env-default:"repository-response"`
	GroupID       string   `yaml:"group_id"`
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

	cfg.GRPCServer.Timeout = time.Duration(cfg.GRPCServer.TimeoutSeconds) * time.Second

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
