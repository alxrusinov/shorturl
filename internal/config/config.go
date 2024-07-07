package config

import "flag"

type Config struct {
	BaseURL     string
	ResponseURL string
}

func (config *Config) Init() {
	flag.StringVar(&config.BaseURL, "a", "localhost:8080", "base url when server will be started")
	flag.StringVar(&config.ResponseURL, "b", "http://localhost:8080", "base url of returning link")

}

func (config *Config) Parse() {
	flag.Parse()
}

func CreateConfig() *Config {
	return &Config{}
}
