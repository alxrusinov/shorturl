package app

import (
	"fmt"

	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/logger"
	"github.com/alxrusinov/shorturl/internal/server"
	"github.com/alxrusinov/shorturl/internal/store"
)

func Run(config *config.Config) {
	store := store.CreateStore(config)
	handler := handler.CreateHandler(store, config.ResponseURL)
	logger := logger.CreateLogger()
	newServer := server.CreateServer(handler, config.BaseURL, logger)

	newServer.Run()

}
