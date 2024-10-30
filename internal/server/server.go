package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/alxrusinov/shorturl/internal/handler"
)

// Server has information about server-mux, handler and server run address
type Server struct {
	mux     *gin.Engine
	handler *handler.Handler
	addr    string
}

// Run runs the server
func (server *Server) Run() {
	server.mux.Run(server.addr)
}

// NewServer initialize and return new server instance
func NewServer(handler *handler.Handler, addr string, logger zerolog.Logger) *Server {
	server := &Server{
		mux:     gin.New(),
		handler: handler,
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

	return server
}
