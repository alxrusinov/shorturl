package config

import (
	"flag"
	"os"
	"sync"
)

// Default fields for config
const (
	// DeafaultBaseURL - base url when server will be started
	DefaulServerAddress = "localhost:8080"
	// DeafaultResponseURL - base url of returning link
	DefaultBaseURL = "http://localhost:8080"
	// DefaultFilePath - path for storage file
	DefaultFilePath = "./config.json"
)

// Config has information about configuration of app
type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DBPath          string `json:"database_dsn"`
	TLS             bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
	ConfigPath      string
}

var once sync.Once

// Init parses flags and initial config
func (config *Config) Init() {
	once.Do(func() {
		flag.StringVar(&config.ServerAddress, "a", DefaulServerAddress, "base url when server will be started")
		flag.StringVar(&config.BaseURL, "b", DefaultBaseURL, "base url of returning link")
		flag.StringVar(&config.FileStoragePath, "f", DefaultFilePath, "path for storage file")
		flag.StringVar(&config.DBPath, "d", "", "path to data base")
		flag.BoolVar(&config.TLS, "s", false, "configure http or https server")
		flag.StringVar(&config.ConfigPath, "c", "", "path to config file")
		flag.StringVar(&config.ConfigPath, "config", "", "path to config file")
		flag.StringVar(&config.TrustedSubnet, "t", "", "trust subnet")
	})

	flag.Parse()

	if baseURL, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.ServerAddress = baseURL
	}

	if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.BaseURL = baseURL
	}

	if filePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = filePath
	}

	if dBPath, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DBPath = dBPath
	}
	if trustedSubnet, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		config.TrustedSubnet = trustedSubnet
	}
	if TLS, ok := os.LookupEnv("ENABLE_HTTPS"); ok && TLS != "" {
		config.TLS = true
	}
	if configPath, ok := os.LookupEnv("CONFIG"); ok {
		config.ConfigPath = configPath
	}

	useConfigFile(config)
}

// NewConfig return Config instance
func NewConfig() *Config {
	return &Config{}
}
