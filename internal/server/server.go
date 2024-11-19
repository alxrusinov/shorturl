package server

import (
	"context"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/alxrusinov/shorturl/internal/config"
	"github.com/alxrusinov/shorturl/internal/handler"
)

// Server has information about server-mux, handler and server run address
type Server struct {
	mux     *gin.Engine
	handler *handler.Handler
	addr    string
	server  *http.Server
	TLS     bool
}

// Run runs the server
func (server *Server) Run() error {
	if server.TLS {
		return server.server.ListenAndServe()
	}

	return server.server.ListenAndServe()

}

// Shutsown realize gracefull shutdown server
func (server *Server) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	if err := server.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// NewServer initialize and return new server instance
func NewServer(handler *handler.Handler, config *config.Config, logger zerolog.Logger) *Server {
	server := &Server{
		mux:     gin.New(),
		handler: handler,
		addr:    config.ServerAddress,
		TLS:     config.TLS,
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

	server.server = &http.Server{
		Addr:    server.addr,
		Handler: server.mux,
	}

	return server
}
