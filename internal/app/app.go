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

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
if (r.Method != http.MethodPost) {
	http.Error(w, "Method is not allowed", http.StatusNotFound)
	return
}


	body, _ := io.ReadAll(r.Body)
	key := string(body)

	shortenURL := generator.GenerateRandomString(10)
	cache[key] = shortenURL

	defer r.Body.Close()

	resp := []byte(shortenURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)


})

mux.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		http.Error(w, "Method is not allowed", http.StatusNotFound)
		return
	}
	id := r.PathValue("id")

	fullURL, ok := cache[id]

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	fmt.Println(id)

	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Location", fullURL)
	w.Write([]byte(""))

})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
