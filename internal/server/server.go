package server

import (
	"log"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/store"
)

type Server struct {
	store store.Store
	mux   *http.ServeMux
}

func (server *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", server.mux))
}

func CreateServer(store store.Store) *Server {
	server := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}

	server.mux.HandleFunc("POST /", GetShortLink(server.store))

	server.mux.HandleFunc("GET /{id}", GetOriginalLink(server.store))

	return server
}
