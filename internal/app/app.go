package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/grpcserver"
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/logger"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/server"
	"github.com/alxrusinov/shorturl/internal/store/dbstore"
	"github.com/alxrusinov/shorturl/internal/store/filestore"
	"github.com/alxrusinov/shorturl/internal/store/inmemorystore"
)

// Run configurate and run application
func Run(ctx context.Context, config *config.Config) {
	var sStore handler.Store

	switch {
	case config.DBPath != "":
		sStore = dbstore.NewDBStore(config.DBPath)
	case config.FileStoragePath != "":
		sStore = filestore.NewFileStore(config.FileStoragePath)
	default:
		sStore = inmemorystore.NewInMemoryStore()

	}

	generator := generator.NewGenerator()

	handler := handler.NewHandler(sStore, config.BaseURL, generator, config.TrustedSubnet)
	logger := logger.NewLogger()
	newServer := server.NewServer(handler, config, logger)
	newGRPSServer := grpcserver.NewGRPCServer(sStore, config.GRPCAddress, generator, config.GRPCAddress, config.TrustedSubnet)

	go func() {
		var batch [][]model.StoreRecord

		for val := range handler.DeleteChan {
			batch = append(batch, val)
			sStore.DeleteLinks(batch)

			batch = batch[0:0]
		}
	}()

	go func() {
		var batch [][]model.StoreRecord

		for val := range handler.DeleteChan {
			batch = append(batch, val)
			sStore.DeleteLinks(batch)

			batch = batch[0:0]
		}
	}()

	go func(ctx context.Context) {
		<-ctx.Done()
		if err := newServer.Shutdown(ctx); err != nil {
			log.Fatal("server has been crashed shutdown")
		}
	}(ctx)

	go func() {
		var batch [][]model.StoreRecord

		for val := range newGRPSServer.DeleteChan {
			batch = append(batch, val)
			sStore.DeleteLinks(batch)

			batch = batch[0:0]
		}
	}()

	if err := newServer.Run(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server has been crashed run")
	}

	if err := grpcserver.Run(newGRPSServer); err != nil {
		log.Fatal("grpc server has been crashed run")
	}

}
