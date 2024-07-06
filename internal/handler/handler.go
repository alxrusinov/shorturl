package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/store"
)

type Handler struct {
	store store.Store
}

func (handler *Handler) GetShortLink(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	originURL := string(body)

	shortenURL := generator.GenerateRandomString(10)
	handler.store.SetLink(shortenURL, originURL)

	defer r.Body.Close()

	resp := []byte(fmt.Sprintf("http://%s/%s", r.Host, shortenURL))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

func (handler *Handler) GetOriginalLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fullURL, err := handler.store.GetLink(id)

	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(""))

}

func CreateHandler(store store.Store) *Handler {
	handler := &Handler{
		store: store,
	}

	return handler
}
