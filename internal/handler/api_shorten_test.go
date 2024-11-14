package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alxrusinov/shorturl/internal/customerrors"
	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/inmemorystore"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func BenchmarkAPIShorten(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testGenerator := mockgenerator.NewMockGenerator()
	testStore := inmemorystore.NewInMemoryStore()
	testHandler := NewHandler(testStore, "http://localhost:8080", testGenerator)
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

func TestHandler_APIShorten(t *testing.T) {
	const userID string = "FOO"

	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator)

	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/api/shorten", testHandler.APIShorten)

	routerWithoutCookie := gin.New()

	routerWithoutCookie.POST("/api/shorten", testHandler.APIShorten)

	testGenerator.On("GenerateRandomString").Return("123", nil)

	cookie := http.Cookie{
		Name:  UserCookie,
		Value: userID,
	}

	tests := []struct {
		name      string
		code      int
		result    APIShortenResult
		body      APIShortenBody
		err       error
		wrongBody string
	}{
		{
			name:   "#1 success",
			code:   http.StatusCreated,
			body:   APIShortenBody{URL: "http://example.com"},
			result: APIShortenResult{Result: "http://localhost:8080/123"},
		},
		{
			name:   "#2 duplicate error",
			code:   http.StatusConflict,
			body:   APIShortenBody{URL: "http://example.com"},
			result: APIShortenResult{Result: "http://localhost:8080/123"},
		},
		{
			name:   "#3 another error",
			code:   http.StatusInternalServerError,
			body:   APIShortenBody{URL: "http://example.com"},
			result: APIShortenResult{Result: "http://localhost:8080/123"},
		},
		{
			name:      "#4 wrong body",
			code:      http.StatusInternalServerError,
			body:      APIShortenBody{URL: "http://example.com"},
			result:    APIShortenResult{Result: "http://localhost:8080/123"},
			wrongBody: "12",
		},
		{
			name:   "#5 wrong random",
			code:   http.StatusInternalServerError,
			body:   APIShortenBody{URL: "http://example.com"},
			result: APIShortenResult{Result: "http://localhost:8080/123"},
		},
		{
			name:   "#6 no cookie",
			code:   http.StatusInternalServerError,
			body:   APIShortenBody{URL: "http://example.com"},
			result: APIShortenResult{Result: "http://localhost:8080/123"},
		},
		{
			name:   "#7 wrong type of body",
			code:   http.StatusInternalServerError,
			body:   APIShortenBody{URL: "http://example.com"},
			result: APIShortenResult{Result: "http://localhost:8080/123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var testRouter *gin.Engine

			testRouter = router

			if tt.name == tests[5].name {
				testRouter = routerWithoutCookie
			}

			body, _ := json.Marshal(tt.body)

			if tt.name == tests[3].name {
				body = []byte(tt.wrongBody)
			}

			if tt.name == tests[4].name {
				testGenerator.On("GenerateRandomString").Unset()
				testGenerator.On("GenerateRandomString").Return("", errors.New("random error"))
			}

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewReader(body))

			if tt.name == tests[6].name {
				request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", newErrReader())
			}

			w := httptest.NewRecorder()

			request.AddCookie(&cookie)

			teststore.On("SetLink", mock.Anything).Unset()

			switch tt.name {
			case tests[0].name:
				teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{
					ShortLink:    "123",
					OriginalLink: tt.body.URL,
					UUID:         userID,
				}, nil)
			case tests[1].name:
				teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{
					ShortLink:    "123",
					OriginalLink: tt.body.URL,
					UUID:         userID,
				}, &customerrors.DuplicateValueError{
					Err: errors.New("duplicate error"),
				})
			case tests[2].name:
				teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{
					ShortLink:    "123",
					OriginalLink: tt.body.URL,
					UUID:         userID,
				}, errors.New("duplicate error"))
			}

			testRouter.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)

		})
	}
}
