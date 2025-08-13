package config

import (
	"flag"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP       HTTPConfig       `yaml:"HTTP"`
	PriceCache PriceCacheConfig `yaml:"PriceCache"`
}

type HTTPConfig struct {
	Port int `yaml:"port" env:"HTTP_PORT"`
}

type PriceCacheConfig struct {
	Host string `yaml:"host" env:"PRICE_CACHE_HOST"`
	Port int    `yaml:"port" env:"PRICE_CACHE_PORT"`
}

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		path := FetchConfigPath()
		instance = LoadConfigByPath(path)
	})
	return instance
}

func FetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config-path", "", "config file path")
	flag.Parse()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	if path == "" {
		path = "./config/config.yaml"
	}
	return path
}

func LoadConfigByPath(path string) *Config {
	var cfg Config
	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic(err)
		}
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
