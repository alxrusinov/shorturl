package handler

import (
	"net/http"
	"net/http/httptest"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
)

func ExampleHandler_Ping() {
	store := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()

	handler := NewHandler(store, "http://example.com:8080", testGenerator, "")

	router := gin.New()

	router.GET("/ping", handler.Ping)

	request := httptest.NewRequest(http.MethodGet, "http://example.com:8080/ping", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, request)

	res := w.Result()

	defer res.Body.Close()

}
