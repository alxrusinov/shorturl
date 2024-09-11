package app

import (
	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/logger"
	"github.com/alxrusinov/shorturl/internal/server"
	"github.com/alxrusinov/shorturl/internal/store"
)

func Run(config *config.Config) {
	sStore := store.CreateStore(config)
	handler := handler.CreateHandler(sStore, config.ResponseURL)
	logger := logger.CreateLogger()
	newServer := server.CreateServer(handler, config.BaseURL, logger)

	go func() {
		var batch [][]store.StoreRecord

		for val := range handler.DeleteChan {
			batch = append(batch, val)
			sStore.DeleteLinks(batch)

			batch = batch[0:0]
		}
	}()

	newServer.Run()

}
