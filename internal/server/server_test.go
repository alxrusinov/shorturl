package server

import (
	"reflect"
	"testing"

	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/logger"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
	}{
		{name: "1# success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := "http://ex.com"
			testStore := mockstore.NewMockStore()
			testGenerator := mockgenerator.NewMockGenerator()
			testHandler := handler.NewHandler(testStore, addr, testGenerator, "")
			logger := logger.NewLogger()

			server := &Server{
				mux:     gin.New(),
				handler: testHandler,
				addr:    addr,
			}

			server.mux.Use(server.handler.Middlewares.LoggerMiddleware(logger))

			server.mux.Use(server.handler.Middlewares.CompressMiddleware())

			server.mux.Use(server.handler.Middlewares.CookieMiddleware())

			server.mux.POST("/", server.handler.GetShortLink)

			server.mux.GET("/:id", server.handler.GetOriginalLink)

			server.mux.POST("/api/shorten", server.handler.APIShorten)

			server.mux.GET("/ping", server.handler.Ping)

			server.mux.POST("/api/shorten/batch", server.handler.APIShortenBatch)

			server.mux.GET("/api/user/urls", server.handler.GetUserLinks)

			server.mux.DELETE("/api/user/urls", server.handler.APIDeleteLinks)

			pprof.Register(server.mux)

			config := config.NewConfig()

			config.ServerAddress = addr

			got := NewServer(testHandler, config, logger)

			if !reflect.DeepEqual(got.handler, server.handler) {
				t.Errorf("NewServer() = %v, want %v", got, server)
			}
			assert.Equal(t, got.addr, server.addr)
		})
	}
}
