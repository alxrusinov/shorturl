package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/store"
	"github.com/gin-gonic/gin"
)

type options struct {
	responseAddr string
}

type Handler struct {
	store   store.Store
	options *options
}

func (handler *Handler) GetShortLink(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	originURL := string(body)

	shortenURL := generator.GenerateRandomString(10)
	handler.store.SetLink(shortenURL, originURL)

	defer ctx.Request.Body.Close()

	resp := []byte(fmt.Sprintf("%s/%s", handler.options.responseAddr, shortenURL))

	ctx.Data(http.StatusCreated, "text/plain", resp)
}

func (handler *Handler) GetOriginalLink(ctx *gin.Context) {
	id := ctx.Param("id")
	defer ctx.Request.Body.Close()

	fullURL, err := handler.store.GetLink(id)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Header("Location", fullURL)
	ctx.Status(http.StatusTemporaryRedirect)
}

func (handler *Handler) APIShorten(ctx *gin.Context) {
	content := struct {
		URL string `json:"url"`
	}{}

	result := struct {
		Result string `json:"result"`
	}{}

	var shortenURL string

	if err := json.NewDecoder(ctx.Request.Body).Decode(&content); err != nil && err != io.EOF {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	// link, err := handler.store.GetLink(content.URL)

	// if err != nil {
	shortenURL = generator.GenerateRandomString(10)
	handler.store.SetLink(shortenURL, content.URL)
	// } else {
	// shortenURL = link
	// }

	result.Result = fmt.Sprintf("%s/%s", handler.options.responseAddr, shortenURL)

	resp, err := json.Marshal(&result)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Data(http.StatusCreated, "application/json", resp)

}

func CreateHandler(store store.Store, responseAddr string) *Handler {
	handler := &Handler{
		store: store,
		options: &options{
			responseAddr: responseAddr,
		},
	}

	return handler
}
