package app

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/generator"
)

func Run() {
	mux := http.NewServeMux()

	cache := make(map[string]string)

mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	originURL := string(body)

	shortenURL := fmt.Sprintf("http://%s/%s",r.Host, generator.GenerateRandomString(10))
	cache[shortenURL] = originURL

	defer r.Body.Close()

	resp := []byte(shortenURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)


})

mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fullURL, ok := cache[id]

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(""))

})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
