package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Init(t *testing.T) {

	config := &Config{}

	tests := []struct {
		name            string
		config          *Config
		osArgs          []string
		wantBaseURL     string
		wantDBPath      string
		wantFilePath    string
		wantResponseURL string
		env             map[string]string
	}{
		{
			name:            "1# base url from args",
			config:          config,
			osArgs:          []string{"test", "-a", "baseURL"},
			wantBaseURL:     "baseURL",
			wantDBPath:      "",
			wantFilePath:    DefaultFilePath,
			wantResponseURL: DeafaultResponseURL,
			env:             map[string]string{},
		},
		{
			name:            "2# db url from args",
			config:          config,
			osArgs:          []string{"test", "-d", "DB", "-a", DeafaultBaseURL},
			wantBaseURL:     DeafaultBaseURL,
			wantDBPath:      "DB",
			wantFilePath:    DefaultFilePath,
			wantResponseURL: DeafaultResponseURL,
			env:             map[string]string{},
		},
		{
			name:            "3# file url from args",
			config:          config,
			osArgs:          []string{"test", "-d", "", "-a", DeafaultBaseURL, "-f", "filePath"},
			wantBaseURL:     DeafaultBaseURL,
			wantDBPath:      "",
			wantFilePath:    "filePath",
			wantResponseURL: DeafaultResponseURL,
			env:             map[string]string{},
		},
		{
			name:            "4# server address env",
			config:          config,
			osArgs:          []string{"test", "-d", "", "-a", DeafaultBaseURL, "-f", DefaultFilePath, "-b", DeafaultResponseURL},
			wantBaseURL:     "foo",
			wantDBPath:      "",
			wantFilePath:    DefaultFilePath,
			wantResponseURL: DeafaultResponseURL,
			env: map[string]string{
				"SERVER_ADDRESS": "foo",
			},
		},
		{
			name:            "5# resp address env",
			config:          config,
			osArgs:          []string{"test", "-d", "", "-a", DeafaultBaseURL, "-f", DefaultFilePath, "-b", DeafaultResponseURL},
			wantBaseURL:     "foo",
			wantDBPath:      "",
			wantFilePath:    DefaultFilePath,
			wantResponseURL: "bar",
			env: map[string]string{
				"SERVER_ADDRESS": "foo",
				"BASE_URL":       "bar",
			},
		},
		{
			name:            "6# server address env",
			config:          config,
			osArgs:          []string{"test", "-d", "", "-a", DeafaultBaseURL, "-f", DefaultFilePath, "-b", DeafaultResponseURL},
			wantBaseURL:     "foo",
			wantDBPath:      "",
			wantFilePath:    "baz",
			wantResponseURL: "bar",
			env: map[string]string{
				"SERVER_ADDRESS":    "foo",
				"BASE_URL":          "bar",
				"FILE_STORAGE_PATH": "baz",
			},
		},
		{
			name:            "7# db address env",
			config:          config,
			osArgs:          []string{"test", "-d", "", "-a", DeafaultBaseURL, "-f", DefaultFilePath, "-b", DeafaultResponseURL},
			wantBaseURL:     "foo",
			wantDBPath:      "clown",
			wantFilePath:    "baz",
			wantResponseURL: "bar",
			env: map[string]string{
				"SERVER_ADDRESS":    "foo",
				"BASE_URL":          "bar",
				"FILE_STORAGE_PATH": "baz",
				"DATABASE_DSN":      "clown",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.osArgs

			for key := range tt.env {
				os.Unsetenv(key)
			}

			for key, val := range tt.env {
				os.Setenv(key, val)
			}

			tt.config.Init()

			assert.Equal(t, tt.wantBaseURL, tt.config.BaseURL)
			assert.Equal(t, tt.wantDBPath, tt.config.DBPath)
			assert.Equal(t, tt.wantFilePath, tt.config.FileStoragePath)
			assert.Equal(t, tt.wantResponseURL, tt.config.ResponseURL)

		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "1# success",
			want: &Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
