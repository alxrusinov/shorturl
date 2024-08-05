package app

import (
	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/logger"
	"github.com/alxrusinov/shorturl/internal/server"
	"github.com/alxrusinov/shorturl/internal/store"
)

func Run(config *config.Config) {
	store := store.CreateFileStore(config.FileStoragePath)
	handler := handler.CreateHandler(store, config.ResponseURL)
	logger := logger.CreateLogger()
	newServer := server.CreateServer(handler, config.BaseURL, logger)

	newServer.Run()

}
