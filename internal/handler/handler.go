package handler

import (
	"encoding/json"
	"errors"
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

	res, err := handler.store.SetLink(links)

	dbErr := &store.DuplicateValueError{}

	if err != nil {
		if !errors.As(err, &dbErr) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return

		}
		dbErr.Err = err
	}

	links.ShortLink = res.ShortLink

	defer ctx.Request.Body.Close()

	resp := []byte(fmt.Sprintf("%s/%s", handler.options.responseAddr, links.ShortLink))

	fmt.Printf("DB ERR: %#v\n", dbErr)

	if dbErr.Err != nil {
		ctx.Data(http.StatusConflict, "text/plain", resp)
		return
	}

	ctx.Data(http.StatusCreated, "text/plain", resp)
}

func (handler *Handler) GetOriginalLink(ctx *gin.Context) {
	id := ctx.Param("id")
	defer ctx.Request.Body.Close()

	links := &store.StoreArgs{
		ShortLink: id,
	}

	res, err := handler.store.GetLink(links)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Header("Location", res.OriginalLink)
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

	if err := json.NewDecoder(ctx.Request.Body).Decode(&content); err != nil && !errors.Is(err, io.EOF) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	shortenURL = generator.GenerateRandomString(10)

	links := &store.StoreArgs{
		ShortLink:    shortenURL,
		OriginalLink: content.URL,
	}

	res, err := handler.store.SetLink(links)

	dbErr := &store.DuplicateValueError{}

	if err != nil {
		if !errors.As(err, &dbErr) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return

		}

		dbErr.Err = err
	}

	links.ShortLink = res.ShortLink

	result.Result = fmt.Sprintf("%s/%s", handler.options.responseAddr, links.ShortLink)

	resp, err := json.Marshal(&result)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if dbErr.Err != nil {
		ctx.Data(http.StatusConflict, "application/json", resp)
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

func (handler *Handler) APIShortenBatch(ctx *gin.Context) {
	var content []*store.StoreArgs

	if err := json.NewDecoder(ctx.Request.Body).Decode(&content); err != nil && !errors.Is(err, io.EOF) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	for _, val := range content {
		shortenURL := generator.GenerateRandomString(10)
		val.ShortLink = shortenURL
	}

	result, err := handler.store.SetBatchLink(content)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, val := range result {
		val.ShortLink = fmt.Sprintf("%s/%s", handler.options.responseAddr, val.ShortLink)
		val.OriginalLink = ""
	}

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
