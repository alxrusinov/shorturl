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

	links := &store.StoreArgs{
		ShortLink:    shortenURL,
		OriginalLink: originURL,
	}

	if err := handler.store.SetLink(links); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	defer ctx.Request.Body.Close()

	resp := []byte(fmt.Sprintf("%s/%s", handler.options.responseAddr, links.ShortLink))

	ctx.Data(http.StatusCreated, "text/plain", resp)
}

func (handler *Handler) GetOriginalLink(ctx *gin.Context) {
	id := ctx.Param("id")
	defer ctx.Request.Body.Close()

	links := &store.StoreArgs{
		ShortLink: id,
	}

	originalURL, err := handler.store.GetLink(links)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Header("Location", originalURL)
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

	shortenURL = generator.GenerateRandomString(10)

	links := &store.StoreArgs{
		ShortLink:    shortenURL,
		OriginalLink: content.URL,
	}

	if err := handler.store.SetLink(links); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	result.Result = fmt.Sprintf("%s/%s", handler.options.responseAddr, links.ShortLink)

	resp, err := json.Marshal(&result)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Data(http.StatusCreated, "application/json", resp)

}

func (handler *Handler) Ping(ctx *gin.Context) {
	err := handler.store.Ping()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)

}

func (handler *Handler) APIShortenBatch(ctx *gin.Context) {}

func CreateHandler(store store.Store, responseAddr string) *Handler {
	handler := &Handler{
		store: store,
		options: &options{
			responseAddr: responseAddr,
		},
	}

	return handler
}
