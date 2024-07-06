package app

import (
	"github.com/alxrusinov/shorturl/internal/server"
	"github.com/alxrusinov/shorturl/internal/store"
)

func Run() {
	store := store.CreateStore()
	newServer := server.CreateServer(store)

	newServer.Run()

}
