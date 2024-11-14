package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_APIDeleteLinks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teststore := new(mockstore.MockStore)
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator)

	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())
	router.DELETE("/api/user/urls", testHandler.APIDeleteLinks)

	routerWithoutCookie := gin.New()

	routerWithoutCookie.DELETE("/api/user/urls", testHandler.APIDeleteLinks)

	tests := []struct {
		name      string
		links     []string
		code      int
		cookie    http.Cookie
		wrongBody map[string]string
	}{
		{
			name:  "1# no user id",
			links: []string{"foo", "bar", "baz"},
			code:  http.StatusInternalServerError,
			cookie: http.Cookie{
				Name:     UserCookie,
				Value:    "fooooooooooo",
				Path:     "/",
				Domain:   "localhost",
				Secure:   false,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour),
			},
		},
		{
			name:  "2# success",
			links: []string{"foo", "bar", "baz"},
			code:  http.StatusAccepted,
			cookie: http.Cookie{
				Name:     UserCookie,
				Value:    "fooooooooooo",
				Path:     "/",
				Domain:   "localhost",
				Secure:   false,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour),
			},
		},
		{
			name:  "3# wrong body",
			links: []string{"foo", "bar", "baz"},
			code:  http.StatusNotFound,
			cookie: http.Cookie{
				Name:     UserCookie,
				Value:    "fooooooooooo",
				Path:     "/",
				Domain:   "localhost",
				Secure:   false,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour),
			},
			wrongBody: map[string]string{},
		},
	}

	for _, test := range tests {
		teststore.On("DeleteLinks", mock.Anything).Unset()
		t.Run(test.name, func(t *testing.T) {
			var testRouter *gin.Engine

			teststore.On("DeleteLinks", mock.Anything).Return(nil)

			testRouter = router

			if test.name == tests[0].name {
				testRouter = routerWithoutCookie
			}

			go func() {
				teststore.On("DeleteLinks", mock.Anything).Return(nil)
				var batch [][]model.StoreRecord

				for val := range testHandler.DeleteChan {
					batch = append(batch, val)
					teststore.DeleteLinks(batch)

					batch = batch[0:0]
				}
			}()

			go func() {
				teststore.On("DeleteLinks", mock.Anything).Return(nil)
				var batch [][]model.StoreRecord

				for val := range testHandler.DeleteChan {
					batch = append(batch, val)
					teststore.DeleteLinks(batch)

					batch = batch[0:0]
				}
			}()

			body, _ := json.Marshal(test.links)

			if test.name == tests[2].name {
				body, _ = json.Marshal(test.wrongBody)
			}

			request := httptest.NewRequest(http.MethodDelete, "http://8080/api/user/urls", bytes.NewReader(body))

			w := httptest.NewRecorder()
			request.AddCookie(&test.cookie)

			testRouter.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, test.code, res.StatusCode)

		})
	}

}
