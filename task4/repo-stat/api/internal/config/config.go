package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string   `yaml:"log_level" env-default:"INFO"`
	Services Services `yaml:"services"`
	HTTP     HTTP
}

type Services struct {
	Processor  string `yaml:"processor" env-required:"true"`
	Subscriber string `yaml:"subscriber" env-required:"true"`
}

type HTTP struct {
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type httpRaw struct {
	Port               string `yaml:"port" env-default:":8080"`
	TimeoutSeconds     int    `yaml:"timeout_seconds" env-default:"5"`
	IdleTimeoutSeconds int    `yaml:"idle_timeout_seconds" env-default:"30"`
}

func Load() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	var rawCfg struct {
		LogLevel string   `yaml:"log_level"`
		Services Services `yaml:"services"`
		HTTP     httpRaw  `yaml:"http"`
	}

	if err := cleanenv.ReadConfig(path, &rawCfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &Config{
		LogLevel: rawCfg.LogLevel,
		Services: rawCfg.Services,
		HTTP: HTTP{
			Port:        rawCfg.HTTP.Port,
			Timeout:     time.Duration(rawCfg.HTTP.TimeoutSeconds) * time.Second,
			IdleTimeout: time.Duration(rawCfg.HTTP.IdleTimeoutSeconds) * time.Second,
		},
	}

	return cfg, nil
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
