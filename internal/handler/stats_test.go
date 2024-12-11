package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Stats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator, "176.14.64.0/18")

	testRouter := gin.New()

	testRouter.GET("/api/internal/stats", testHandler.Stats)

	tests := []struct {
		name    string
		xRealIp string
		code    int
		watnErr bool
	}{
		{
			name:    "1# success",
			xRealIp: "176.14.86.83",
			code:    http.StatusOK,
		},
		{
			name:    "2# wron real ip",
			xRealIp: "",
			code:    http.StatusForbidden,
		},
		{
			name:    "3# error",
			xRealIp: "176.14.86.83",
			code:    http.StatusInternalServerError,
			watnErr: true,
		},
	}

	teststore.On("GetStat").Return(&model.StatResponse{URLS: 10, Users: 10}, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.watnErr {
				teststore.On("GetStat").Unset()
				teststore.On("GetStat").Return(&model.StatResponse{URLS: 10, Users: 10}, errors.New("err"))
			}

			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/internal/stats", nil)

			w := httptest.NewRecorder()

			request.Header.Add("X-Real-IP", tt.xRealIp)

			testRouter.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
