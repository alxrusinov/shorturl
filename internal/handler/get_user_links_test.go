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

func TestHandler_GetUserLinks(t *testing.T) {

	const trueUserID string = "1"

	const errUserID string = "2"

	const noContentUserID string = "3"

	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator)

	teststore.On("GetLinks", trueUserID).Return([]model.StoreRecord{
		{
			ShortLink: "123", OriginalLink: "http://example.com",
		},
	}, nil)

	teststore.On("GetLinks", noContentUserID).Return([]model.StoreRecord{}, nil)

	teststore.On("GetLinks", errUserID).Return([]model.StoreRecord{
		{
			ShortLink: "123", OriginalLink: "http://example.com",
		},
	}, errors.New("error"))

	router := gin.New()

	routerWithoutCookie := gin.New()

	routerWithoutCookie.GET("/api/user/urls", testHandler.GetUserLinks)

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.GET("/api/user/urls", testHandler.GetUserLinks)

	tests := []struct {
		id   string
		name string
		code int
		resp []struct {
			Short    string `json:"short_url"`
			Original string `json:"original_url"`
		}
		err error
	}{
		{
			name: "1# success",
			code: http.StatusOK,
			resp: []struct {
				Short    string `json:"short_url"`
				Original string `json:"original_url"`
			}{
				{
					Short:    "http://localhost:8080/123",
					Original: "http://example",
				},
			},
			id: trueUserID,
		},
		{
			name: "2# error",
			code: http.StatusInternalServerError,
			resp: []struct {
				Short    string `json:"short_url"`
				Original string `json:"original_url"`
			}{
				{
					Short:    "http://localhost:8080/123",
					Original: "http://example",
				},
			},
			id: errUserID,
		},
		{
			name: "3# no content",
			code: http.StatusNoContent,
			resp: []struct {
				Short    string `json:"short_url"`
				Original string `json:"original_url"`
			}{
				{
					Short:    "http://localhost:8080/123",
					Original: "http://example",
				},
			},
			id: noContentUserID,
		},
		{
			name: "4# no cookie",
			code: http.StatusInternalServerError,
			resp: []struct {
				Short    string `json:"short_url"`
				Original string `json:"original_url"`
			}{
				{
					Short:    "http://localhost:8080/123",
					Original: "http://example",
				},
			},
			id: noContentUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var testRouter *gin.Engine

			testRouter = router

			if tt.name == tests[3].name {
				testRouter = routerWithoutCookie
			}

			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/urls", nil)

			w := httptest.NewRecorder()

			cookie := http.Cookie{
				Name:  UserCookie,
				Value: tt.id,
			}

			request.AddCookie(&cookie)

			testRouter.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.code, res.StatusCode)

		})
	}
}
