package server

import (
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	mux     *gin.Engine
	handler *handler.Handler
	addr    string
}

func (server *Server) Run() {
	server.mux.Run(server.addr)
}

func CreateServer(handler *handler.Handler, addr string) *Server {
	server := &Server{
		mux:     gin.Default(),
		handler: handler,
		addr:    addr,
	}

	server.mux.POST("/", server.handler.GetShortLink)

	server.mux.GET("/:id", server.handler.GetOriginalLink)

	return server
}
