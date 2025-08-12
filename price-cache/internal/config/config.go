package config

import (
	"flag"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Kafka KafkaConfig `yaml:"KafkaConfig"`
	Redis RedisConfig `yaml:"RedisConfig"`
	GRPC  GRPCConfig  `yaml:"GRPCConfig"`
}

type KafkaConfig struct {
	Addresses []string `yaml:"addresses" env:"ADDRESSES" env-separator:","`
	Topic     string   `yaml:"topic" env:"TOPIC"`
	GroupID   string   `yaml:"groupId" env:"GROUP_ID"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env:"GRPC_PORT"`
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
