package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/logger"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func Test_checkContentType(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1# true",
			args: args{
				values: []string{"text/html", "application/json"},
			},
			want: true,
		},
		{
			name: "1# false",
			args: args{
				values: []string{"text", "json"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkContentType(tt.args.values); got != tt.want {
				t.Errorf("checkContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkGzip(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1# true",
			args: args{
				values: []string{"gzip", "json"},
			},
			want: true,
		},
		{
			name: "1# false",
			args: args{
				values: []string{"text", "json"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkGzip(tt.args.values); got != tt.want {
				t.Errorf("checkGzip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiddlewares_CompressMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator, "")

	testGenerator.On("GenerateUserID").Return("123", nil)
	testGenerator.On("GenerateRandomString").Return("321", nil)
	teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "short", OriginalLink: "original", UUID: "1"}, nil)

	router := gin.New()

	router.Use(testHandler.Middlewares.CompressMiddleware())
	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/", testHandler.GetShortLink)

	tests := []struct {
		name   string
		reqZip bool
		resZip bool
	}{
		{
			name:   "1# without zip",
			reqZip: false,
			resZip: false,
		},
		{
			name:   "2# all zip",
			reqZip: true,
			resZip: true,
		},
		{
			name:   "3# req zip",
			reqZip: true,
			resZip: false,
		},
		{
			name:   "4# res zip",
			reqZip: false,
			resZip: true,
		},
		{
			name:   "5# wrong content-type",
			reqZip: false,
			resZip: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal("link")

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(body))

			if tt.reqZip {
				request.Header.Add("Content-Encoding", "gzip")
			}

			w := httptest.NewRecorder()

			if tt.resZip {
				w.Header().Add("Content-Encoding", "gzip")
				w.Header().Add("Accept-Encoding", "gzip")

			}

			if tt.name == tests[4].name {
				request.Method = http.MethodGet
				request.Header.Set("Content-Type", "application/pdf")
			}

			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

		})
	}
}

func TestMiddlewares_CookieMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator, "")

	testGenerator.On("GenerateUserID").Return("", errors.New("err"))
	testGenerator.On("GenerateRandomString").Return("321", nil)
	teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "short", OriginalLink: "original", UUID: "1"}, nil)

	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/", testHandler.GetShortLink)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "1# user id error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal("link")

			if tt.wantErr {
				testGenerator.On("GenerateUserID").Unset()
				testGenerator.On("GenerateUserID").Return("", errors.New("err"))
			}

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(body))

			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

		})
	}
}

func TestMiddlewares_LoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator, "")
	logger := logger.NewLogger()

	testGenerator.On("GenerateUserID").Return("111", nil)
	testGenerator.On("GenerateRandomString").Return("321", nil)
	teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "short", OriginalLink: "original", UUID: "1"}, nil)

	router := gin.New()

	router.Use(testHandler.Middlewares.LoggerMiddleware(logger))

	router.POST("/", testHandler.GetShortLink)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "1# success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal("link")

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(body))

			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

		})
	}
}
