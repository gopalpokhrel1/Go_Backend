package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var config_path string

	config_path = os.Getenv("CONFIG_PATH")

	if config_path == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		config_path = *flags

		fmt.Println(config_path)

		if config_path == "" {
			log.Fatal("config path is not set")
		}
	}

	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Fatal("Config file is not exists")
	}

	var cfg Config
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		log.Fatalf("Cannot read config file %v", err)
	}

	return &cfg
}
