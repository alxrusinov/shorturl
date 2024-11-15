package handler

import (
	"bytes"
	"encoding/json"
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

func TestHandler_APIShortenBatch(t *testing.T) {
	const userID string = "FOO"

	gin.SetMode(gin.TestMode)
	teststore := new(mockstore.MockStore)
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator)

	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/api/shorten/batch", testHandler.APIShortenBatch)

	routerWithoutCookie := gin.New()

	routerWithoutCookie.POST("/api/shorten/batch", testHandler.APIShortenBatch)

	testGenerator.On("GenerateRandomString").Return("123", nil)

	cookie := http.Cookie{
		Name:  UserCookie,
		Value: userID,
	}

	tests := []struct {
		name      string
		code      int
		result    []*model.StoreRecord
		body      []*model.StoreRecord
		err       error
		wrongBody string
	}{
		{
			name: "1# success",
			code: http.StatusCreated,
			body: []*model.StoreRecord{
				{
					CorrelationID: "1",
					OriginalLink:  "http://clown.com",
				},
				{
					CorrelationID: "2",
					OriginalLink:  `http://example.com`,
				},
			},
			result: []*model.StoreRecord{
				{
					CorrelationID: "1",
					ShortLink:     "111",
				},
				{
					CorrelationID: "2",
					ShortLink:     `222`,
				},
			},
		},
		{
			name: "2# error",
			code: http.StatusInternalServerError,
			body: []*model.StoreRecord{
				{
					CorrelationID: "1",
					OriginalLink:  "http://clown.com",
				},
				{
					CorrelationID: "2",
					OriginalLink:  `http://example.com`,
				},
			},
			result: []*model.StoreRecord{
				{
					CorrelationID: "1",
					ShortLink:     "111",
				},
				{
					CorrelationID: "2",
					ShortLink:     `222`,
				},
			},
			err: errors.New("error"),
		},
		{
			name: "3# wrong body",
			code: http.StatusNotFound,
			body: []*model.StoreRecord{
				{
					CorrelationID: "1",
					OriginalLink:  "http://clown.com",
				},
				{
					CorrelationID: "2",
					OriginalLink:  `http://example.com`,
				},
			},
			result: []*model.StoreRecord{
				{
					CorrelationID: "1",
					ShortLink:     "111",
				},
				{
					CorrelationID: "2",
					ShortLink:     `222`,
				},
			},
			wrongBody: "foo",
		},
		{
			name: "4# wrong random",
			code: http.StatusInternalServerError,
			body: []*model.StoreRecord{
				{
					CorrelationID: "1",
					OriginalLink:  "http://clown.com",
				},
				{
					CorrelationID: "2",
					OriginalLink:  `http://example.com`,
				},
			},
			result: []*model.StoreRecord{
				{
					CorrelationID: "1",
					ShortLink:     "111",
				},
				{
					CorrelationID: "2",
					ShortLink:     `222`,
				},
			},
		},
		{
			name: "5# no cookie",
			code: http.StatusInternalServerError,
			body: []*model.StoreRecord{
				{
					CorrelationID: "1",
					OriginalLink:  "http://clown.com",
				},
				{
					CorrelationID: "2",
					OriginalLink:  `http://example.com`,
				},
			},
			result: []*model.StoreRecord{
				{
					CorrelationID: "1",
					ShortLink:     "111",
				},
				{
					CorrelationID: "2",
					ShortLink:     `222`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var testRouter *gin.Engine

			testRouter = router
			body, _ := json.Marshal(tt.body)

			if tt.name == tests[2].name {
				body = []byte(tt.wrongBody)
			}

			if tt.name == tests[3].name {
				testGenerator.On("GenerateRandomString").Unset()
				testGenerator.On("GenerateRandomString").Return("", errors.New("random error"))

			}

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/batch", bytes.NewReader(body))

			w := httptest.NewRecorder()

			request.AddCookie(&cookie)

			teststore.On("SetBatchLink", mock.Anything).Unset()

			switch tt.name {
			case tests[0].name:
				teststore.On("SetBatchLink", mock.Anything).Return(tt.result, nil)
			case tests[1].name:
				teststore.On("SetBatchLink", mock.Anything).Return(tt.result, tt.err)
			case tests[2].name:
				teststore.On("SetBatchLink", mock.Anything).Return(tt.result, nil)
			case tests[4].name:
				testRouter = routerWithoutCookie
			}

			testRouter.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
