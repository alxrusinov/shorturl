package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	Host    = "localhost:8080"
	Scheme = "http"
)

func Run() {
	urlAddr := url.URL{
		Scheme: Scheme,
		Host: Host,
	}

	mux := http.NewServeMux()

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
if (r.Method != http.MethodPost) {
	http.Error(w, "Method is not allowed", http.StatusNotFound)
	return
}


	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	fmt.Println(body)

	shortenUrl := "foo"

	resp := []byte(shortenUrl)

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

	fullUrl := "http://full-url.com"

	fmt.Println(id)

	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Location", fullUrl)
	w.Write([]byte(""))

})

	log.Fatal(http.ListenAndServe(urlAddr.String(), mux))
}
