package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// NewLogger creates instance of logger
func NewLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stdout)

	return logger
}
