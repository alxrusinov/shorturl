package app

import (
	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/server"
	"github.com/alxrusinov/shorturl/internal/store"
)

func Run(config *config.Config) {
	store := store.CreateStore()
	handler := handler.CreateHandler(store, config.ResponseURL)
	newServer := server.CreateServer(handler, config.BaseURL)

	newServer.Run()

}
