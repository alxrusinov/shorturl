package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alxrusinov/shorturl/internal/store/inmemorystore"
	"github.com/gin-gonic/gin"
)

func BenchmarkAPIShorten(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testStore := inmemorystore.NewInMemoryStore()
	testHandler := NewHandler(testStore, "http://localhost:8080")
	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/api/shorten", testHandler.APIShorten)

	content := struct {
		URL string `json:"url"`
	}{URL: "http://example.com"}

	send, _ := json.Marshal(&content)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewReader(send))
		w := httptest.NewRecorder()
		b.StartTimer()
		router.ServeHTTP(w, request)

	}
}
