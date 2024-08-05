package config

import (
	"flag"
	"os"
)

const (
	DeafaultBaseURL        = "localhost:8080"
	DeafaultResponseURL    = "http://localhost:8080"
	DefaultFileStoragePath = "./config.txt"
)

type Config struct {
	BaseURL         string
	ResponseURL     string
	FileStoragePath string
}

func (config *Config) Init() {
	flag.StringVar(&config.BaseURL, "a", DeafaultBaseURL, "base url when server will be started")
	flag.StringVar(&config.ResponseURL, "b", DeafaultResponseURL, "base url of returning link")
	flag.StringVar(&config.FileStoragePath, "f", DefaultFileStoragePath, "path for storage file")
}

func (config *Config) Parse() {
	flag.Parse()

	if baseURL, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.BaseURL = baseURL
	}

	if responseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.ResponseURL = responseURL
	}

	if filePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = filePath
	}
}

func CreateConfig() *Config {
	return &Config{}
}
