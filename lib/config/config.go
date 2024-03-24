package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server   `yaml:"server"`
	Db     Database `yaml:"db"`
}

type Server struct {
	Port            int    `yaml:"port"`
	IntegerServAddr string `yaml:"integrserv"`
}

type Database struct {
	DriverName string `yaml:"driver"`
	SourcePath string `yaml:"source"`
}

func MustLoad() *Config {
	path := fetchConfigByPath()
	if path == "" {
		panic("Failed to loading config: Config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	stream, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(stream, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigByPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
