package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`

	Jwt JwtConfig `yaml:"jwt"`
}

const defaultPath = "config/config.yaml"

func LoadFromYAML() *Config {
	filePath := getPathFromCmdArgs()
	if filePath == "" {
		filePath = defaultPath
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exists: %s\n%s", filePath, err)
	}

	config := Config{}
	if err := cleanenv.ReadConfig(filePath, &config); err != nil {
		log.Fatalf("can't read config file: %s", err)
	}

	return &config
}

func getPathFromCmdArgs() string {
	path := ""

	flag.StringVar(&path, "config", "", "config path")
	flag.Parse()

	return path
}
