package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alxrusinov/shorturl/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetOriginalLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testStore := store.CreateStore()
	testHandler := CreateHandler(testStore)
	router := gin.New()

	router.GET("/:id", testHandler.GetOriginalLink)

	originalLink := "http://example.com"
	shortenLink := "abcde"

	testHandler.store.SetLink(shortenLink, originalLink)

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
				response: originalLink,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/abcde", nil)

			fmt.Printf("%+v", request)

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
	testStore := store.CreateStore()
	testHandler := CreateHandler(testStore)
	router := gin.New()

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
