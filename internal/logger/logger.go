package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func CreateLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stdout)

	return logger
}
