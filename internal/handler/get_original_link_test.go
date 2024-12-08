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
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetOriginalLink(t *testing.T) {
	const userID string = "FOO"

	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator, "")

	router := gin.New()

	router.GET("/:id", testHandler.GetOriginalLink)

	cookie := http.Cookie{
		Name:  UserCookie,
		Value: userID,
	}

	tests := []struct {
		name string
		code int
		resp string
		err  error
		id   string
	}{
		{
			name: "1# success",
			code: http.StatusTemporaryRedirect,
			resp: "111",
			err:  nil,
			id:   "123",
		},
		{
			name: "2# error",
			code: http.StatusGone,
			resp: "222",
			err:  errors.New("err"),
			id:   "321",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/"+tt.id, nil)

			teststore.On("GetLink", mock.Anything).Unset()

			switch tt.name {
			case tests[0].name:
				teststore.On("GetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "123", OriginalLink: "111"}, nil)
			case tests[1].name:
				teststore.On("GetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "123", OriginalLink: "222"}, errors.New("err"))
			}

			w := httptest.NewRecorder()

			request.AddCookie(&cookie)

			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)

			if tt.err == nil {
				assert.Equal(t, tt.resp, res.Header.Get("Location"))

			}
		})
	}
}
