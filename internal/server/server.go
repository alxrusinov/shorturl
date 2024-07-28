package server

import (
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	mux     *gin.Engine
	handler *handler.Handler
	addr    string
}

func (server *Server) Run() {
	server.mux.Run(server.addr)
}

func CreateServer(handler *handler.Handler, addr string, logger zerolog.Logger) *Server {
	server := &Server{
		mux:     gin.New(),
		handler: handler,
		addr:    addr,
	}

	server.mux.Use(loggerMiddleware(logger))

	server.mux.POST("/", server.handler.GetShortLink)

	server.mux.GET("/:id", server.handler.GetOriginalLink)

	server.mux.POST("/api/shorten", server.handler.APIShorten)

	return server
}
