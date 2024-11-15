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
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetShortLink(t *testing.T) {
	const userID string = "FOO"
	const successCase string = "successCase"
	const errCase string = "errCase"
	const duplicateErrCase string = "duplicateErrCase"

	gin.SetMode(gin.TestMode)
	teststore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	testHandler := NewHandler(teststore, "http://localhost:8080", testGenerator)

	router := gin.New()

	router.Use(testHandler.Middlewares.CookieMiddleware())

	router.POST("/", testHandler.GetShortLink)

	routerWithoutCookie := gin.New()
	routerWithoutCookie.POST("/", testHandler.GetShortLink)

	testGenerator.On("GenerateRandomString").Return("123", nil)

	cookie := http.Cookie{
		Name:  UserCookie,
		Value: userID,
	}

	tests := []struct {
		name      string
		testcase  string
		link      string
		resp      string
		code      int
		wrongBody map[string]struct{}
	}{
		{
			name:      "#1 success",
			testcase:  successCase,
			link:      "http://example.com",
			resp:      "123",
			code:      http.StatusCreated,
			wrongBody: nil,
		},
		{
			name:      "#2 duplicate value",
			testcase:  duplicateErrCase,
			link:      "http://example.com",
			resp:      "123",
			code:      http.StatusConflict,
			wrongBody: nil,
		},
		{
			name:      "#3 some error",
			testcase:  errCase,
			link:      "http://example.com",
			resp:      "123",
			code:      http.StatusInternalServerError,
			wrongBody: nil,
		},
		{
			name:      "#4 random error",
			testcase:  errCase,
			link:      "http://example.com",
			resp:      "123",
			code:      http.StatusInternalServerError,
			wrongBody: nil,
		},
		{
			name:      "#5 no cookie",
			testcase:  errCase,
			link:      "http://example.com",
			resp:      "123",
			code:      http.StatusInternalServerError,
			wrongBody: nil,
		},
		{
			name:      "#6 no cookie",
			testcase:  errCase,
			link:      "http://example.com",
			resp:      "123",
			code:      http.StatusInternalServerError,
			wrongBody: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testRouter *gin.Engine

			body, _ := json.Marshal(tt.link)

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(body))

			if tt.name == tests[5].name {
				request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/", newErrReader())
			}

			w := httptest.NewRecorder()

			request.AddCookie(&cookie)

			teststore.On("SetLink").Unset()

			testRouter = router

			if tt.name == tests[3].name {
				testGenerator.On("GenerateRandomString").Unset()
				testGenerator.On("GenerateRandomString").Return("", errors.New("random error"))
			}

			if tt.name == tests[4].name {
				testRouter = routerWithoutCookie
			}

			switch tt.testcase {
			case successCase:
				teststore.On("SetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: tt.resp, OriginalLink: tt.link, UUID: userID}, nil)
			case duplicateErrCase:
				teststore.On("SetLink", mock.Anything).Return(new(model.StoreRecord), &customerrors.DuplicateValueError{Err: errors.New("duplicate error")})
			case errCase:
				teststore.On("SetLink", mock.Anything).Return(new(model.StoreRecord), errors.New("error"))
			}

			testRouter.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
