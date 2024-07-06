package server

import (
	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/store"
	"github.com/gin-gonic/gin"
)

type Server struct {
	mux     *gin.Engine
	handler *handler.Handler
}

func (server *Server) Run() {
	server.mux.Run(":8080")
}

func CreateServer(store store.Store) *Server {
	server := &Server{
		mux:     gin.Default(),
		handler: handler.CreateHandler(store),
	}

	server.mux.POST("/", server.handler.GetShortLink)

	server.mux.GET("/:id", server.handler.GetOriginalLink)

	return server
}
