package config

import (
	"flag"
	"os"
)

const (
	DeafaultBaseURL     = "localhost:8080"
	DeafaultResponseURL = "http://localhost:8080"
)

type Config struct {
	BaseURL     string
	ResponseURL string
}

func (config *Config) Init() {
	flag.StringVar(&config.BaseURL, "a", DeafaultBaseURL, "base url when server will be started")
	flag.StringVar(&config.ResponseURL, "b", DeafaultResponseURL, "base url of returning link")

}

func (config *Config) Parse() {
	flag.Parse()

	if baseURL, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.BaseURL = baseURL
	}

	if responseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.ResponseURL = responseURL
	}
}

func CreateConfig() *Config {
	return &Config{}
}
