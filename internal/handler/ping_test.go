package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Ping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teststore := new(mockstore.MockStore)
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator)

	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.GET("/ping", testHandler.Ping)

	tests := []struct {
		name   string
		code   int
		hasErr bool
	}{
		{
			name:   "#1 success",
			code:   http.StatusOK,
			hasErr: false,
		},
		{
			name:   "#2 error",
			code:   http.StatusInternalServerError,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)

			w := httptest.NewRecorder()

			teststore.On("Ping").Unset()

			if tt.hasErr {
				teststore.On("Ping").Return(errors.New("error"))
			} else {
				teststore.On("Ping").Return(nil)
			}

			router.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.code, res.StatusCode)

		})
	}
}
