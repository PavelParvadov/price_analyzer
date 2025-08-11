package config

import (
	"flag"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PriceConfig PriceConfig    `yaml:"PriceConfig"`
	KafkaConfig KafkaConfig    `yaml:"KafkaConfig"`
	Producer    ProducerConfig `yaml:"ProducerConfig"`
}

type PriceConfig struct {
	Port int `yaml:"PricePort" env:"PRICE_PORT"`
}

type KafkaConfig struct {
	Addresses []string `yaml:"addresses" env:"ADDRESSES"`
	Topic     string   `yaml:"topic" env:"TOPIC"`
}
type ProducerConfig struct {
	Tickers           []string `yaml:"tickers" env:"TICKERS" env-separator:","`
	IntervalMs        int      `yaml:"intervalMs" env:"INTERVAL_MS"`
	InitialPrice      float64  `yaml:"initialPrice" env:"INITIAL_PRICE"`
	VolatilityPercent float64  `yaml:"volatilityPercent" env:"VOLATILITY_PERCENT"`
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		path := FetchConfigPath()
		instance = LoadConfigByPath(path)
	})
	return instance
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

func FetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config-path", "", "load config from path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
