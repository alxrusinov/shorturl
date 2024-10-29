package config

import (
	"flag"
	"os"
)

// Default fields for config
const (
	// DeafaultBaseURL - base url when server will be started
	DeafaultBaseURL = "localhost:8080"
	// DeafaultResponseURL - base url of returning link
	DeafaultResponseURL = "http://localhost:8080"
	// DefaultFilePath - path for storage file
	DefaultFilePath = "./config.json"
)

// Config has information about configuration of app
type Config struct {
	BaseURL         string
	ResponseURL     string
	FileStoragePath string
	DBPath          string
}

// Init parses flags and initial config
func (config *Config) Init() {
	flag.StringVar(&config.BaseURL, "a", DeafaultBaseURL, "base url when server will be started")
	flag.StringVar(&config.ResponseURL, "b", DeafaultResponseURL, "base url of returning link")
	flag.StringVar(&config.FileStoragePath, "f", DefaultFilePath, "path for storage file")
	flag.StringVar(&config.DBPath, "d", "", "path to data base")
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

	if dBPath, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DBPath = dBPath
	}
}

// NewConfig return Config instance
func NewConfig() *Config {
	return &Config{}
}
