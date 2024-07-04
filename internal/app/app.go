package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
if (r.Method != http.MethodPost) {
	http.Error(w, "Method is not allowed", http.StatusNotFound)
	return
}


	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	fmt.Println(body)

	shortenURL := "foo"

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

	fullURL := "http://full-url.com"

	fmt.Println(id)

	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Location", fullURL)
	w.Write([]byte(""))

})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
