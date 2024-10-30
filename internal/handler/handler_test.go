package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/inmemorystore"
)

func TestHandler_GetOriginalLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testStore := inmemorystore.NewInMemoryStore()
	testHandler := NewHandler(testStore, "http://localhost:8080")
	router := gin.New()

	router.GET("/:id", testHandler.GetOriginalLink)

	links := &model.StoreRecord{
		ShortLink:    "abcde",
		OriginalLink: "http://example.com",
	}

	testHandler.store.SetLink(links)

	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: links.OriginalLink,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/abcde", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			result := w.Header().Get("Location")

			assert.Equal(t, test.want.response, result)
		})
	}
}

func TestHandler_GetShortLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testStore := inmemorystore.NewInMemoryStore()
	testHandler := NewHandler(testStore, "http://localhost:8080")
	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/", testHandler.GetShortLink)

	type want struct {
		code        int
		response    string
		contentType string
		error       error
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
				response:    "",
				error:       nil,
			},
		},
	}

	for _, test := range tests {
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("http://example.com"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)

		res := w.Result()
		assert.Equal(t, test.want.code, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)

		defer res.Body.Close()
		require.NoError(t, err)

		assert.NotEmpty(t, resBody)
		assert.Equal(t, test.want.contentType, w.Header().Get("Content-Type"))
	}

}

func TestHandler_APIShorten(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testStore := inmemorystore.NewInMemoryStore()
	testHandler := NewHandler(testStore, "http://localhost:8080")
	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/api/shorten", testHandler.APIShorten)

	content := struct {
		URL string `json:"url"`
	}{URL: "http://example.com"}

	result := struct {
		Result string `json:"result"`
	}{}

	type want struct {
		code        int
		response    string
		contentType string
		error       error
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
				response:    "",
				error:       nil,
			},
		},
	}

	send, _ := json.Marshal(&content)

	for _, test := range tests {
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewReader(send))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)

		res := w.Result()

		assert.Equal(t, test.want.code, res.StatusCode)

		if err := json.NewDecoder(res.Body).Decode(&result); err != nil && err != io.EOF {
			require.NoError(t, err)
		}

		defer res.Body.Close()

		assert.NotEmpty(t, result.Result)
		assert.Equal(t, test.want.contentType, w.Header().Get("Content-Type"))
	}

}
