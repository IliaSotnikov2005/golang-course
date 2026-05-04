package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel  string
	Services  Services
	HTTP      HTTP
	Redis     Redis
	Cache     Cache
	RateLimit RateLimit
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

type Redis struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"6379"`
}

type Cache struct {
	TTL time.Duration
}

type cacheRaw struct {
	TTLSeconds int `yaml:"ttl_seconds" env-default:"60"`
}

type RateLimit struct {
	RequestsPerSecond float64 `yaml:"requests_per_second" env-default:"5"`
	Burst             int     `yaml:"burst" env-default:"10"`
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
		LogLevel  string    `yaml:"log_level" env-default:"INFO"`
		Services  Services  `yaml:"services"`
		HTTP      httpRaw   `yaml:"http"`
		Redis     Redis     `yaml:"redis"`
		Cache     cacheRaw  `yaml:"cache"`
		RateLimit RateLimit `yaml:"rate_limit"`
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
		Redis: rawCfg.Redis,
		Cache: Cache{
			TTL: time.Duration(rawCfg.Cache.TTLSeconds) * time.Second,
		},
		RateLimit: rawCfg.RateLimit,
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
