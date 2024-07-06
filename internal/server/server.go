package server

import (
	"log"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/handler"
	"github.com/alxrusinov/shorturl/internal/store"
)

type Server struct {
	mux     *http.ServeMux
	handler *handler.Handler
}

func (server *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", server.mux))
}

func CreateServer(store store.Store) *Server {
	server := &Server{
		mux:     http.NewServeMux(),
		handler: handler.CreateHandler(store),
	}

	server.mux.HandleFunc("POST /", server.handler.GetShortLink)

	server.mux.HandleFunc("GET /{id}", server.handler.GetOriginalLink)

	return server
}
