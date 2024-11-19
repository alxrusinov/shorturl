package config

import (
	"encoding/json"
	"io"
	"os"
)

func readConfig(file *os.File) (*Config, error) {
	content, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	configFromFile := new(Config)

	err = json.Unmarshal(content, &configFromFile)

	if err != nil {
		return nil, err
	}

	return configFromFile, nil
}

func useConfigFile(config *Config) {
	file, err := os.Open(config.ConfigPath)

	if err != nil {
		return
	}

	defer file.Close()

	configFromFile, err := readConfig(file)

	if err != nil {
		return
	}

	if configFromFile.ServerAddress != "" && config.ServerAddress == DefaulServerAddress {
		config.ServerAddress = configFromFile.ServerAddress
	}

	if configFromFile.BaseURL != "" && config.BaseURL == DefaultBaseURL {
		config.BaseURL = configFromFile.BaseURL
	}

	if configFromFile.FileStoragePath != "" && config.FileStoragePath == DefaultFilePath {
		config.FileStoragePath = configFromFile.FileStoragePath
	}

	if configFromFile.DBPath != "" && config.DBPath == "" {
		config.DBPath = configFromFile.DBPath
	}

	if configFromFile.TLS && !config.TLS {
		config.TLS = configFromFile.TLS
	}

}
