package handler

import (
	"encoding/json"
	"errors"
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
	store       store.Store
	options     *options
	Middlewares *Middlewares
	DeleteChan  chan []store.StoreRecord
}

func (handler *Handler) GetShortLink(ctx *gin.Context) {
	var userID string

	val, ok := ctx.Get("userID")

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok = val.(string)

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	body, _ := io.ReadAll(ctx.Request.Body)
	originURL := string(body)

	shortenURL, err := generator.GenerateRandomString()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	links := &store.StoreRecord{
		ShortLink:    shortenURL,
		OriginalLink: originURL,
		UUID:         userID,
	}

	res, err := handler.store.SetLink(links)

	dbErr := &store.DuplicateValueError{}

	if err != nil {
		if !errors.As(err, &dbErr) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return

		}
	}

	links.ShortLink = res.ShortLink

	defer ctx.Request.Body.Close()

	resp := []byte(createShortLink(handler.options.responseAddr, links.ShortLink))

	if dbErr.Err != nil {
		ctx.Data(http.StatusConflict, "text/plain", resp)
		return
	}

	ctx.Data(http.StatusCreated, "text/plain", resp)
}

func (handler *Handler) GetOriginalLink(ctx *gin.Context) {
	id := ctx.Param("id")
	defer ctx.Request.Body.Close()

	links := &store.StoreRecord{
		ShortLink: id,
	}

	res, err := handler.store.GetLink(links)

	if err != nil {
		ctx.Status(http.StatusGone)
		return
	}

	ctx.Header("Location", res.OriginalLink)
	ctx.Status(http.StatusTemporaryRedirect)
}

func (handler *Handler) APIShorten(ctx *gin.Context) {
	var userID string

	val, ok := ctx.Get("userID")

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok = val.(string)

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

	shortenURL, err := generator.GenerateRandomString()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	links := &store.StoreRecord{
		ShortLink:    shortenURL,
		OriginalLink: content.URL,
		UUID:         userID,
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

	result.Result = createShortLink(handler.options.responseAddr, links.ShortLink)

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
	var userID string

	val, ok := ctx.Get("userID")

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok = val.(string)

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var content []*store.StoreRecord

	if err := json.NewDecoder(ctx.Request.Body).Decode(&content); err != nil && !errors.Is(err, io.EOF) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	for _, val := range content {
		shortenURL, err := generator.GenerateRandomString()

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		val.ShortLink = shortenURL
		val.UUID = userID
	}

	result, err := handler.store.SetBatchLink(content)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, val := range result {
		val.ShortLink = createShortLink(handler.options.responseAddr, val.ShortLink)
		val.OriginalLink = ""
	}

	resp, err := json.Marshal(&result)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Data(http.StatusCreated, "application/json", resp)

}

func (handler *Handler) GetUserLinks(ctx *gin.Context) {
	var userID string

	val, ok := ctx.Get("userID")

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok = val.(string)

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	links, err := handler.store.GetLinks(userID)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(links) == 0 {
		ctx.Header("Content-Type", "application/json")
		ctx.Status(http.StatusNoContent)
		return
	}

	var result []struct {
		Short    string `json:"short_url"`
		Original string `json:"original_url"`
	}

	for _, link := range links {
		newLink := struct {
			Short    string `json:"short_url"`
			Original string `json:"original_url"`
		}{
			Short:    createShortLink(handler.options.responseAddr, link.ShortLink),
			Original: link.OriginalLink,
		}
		result = append(result, newLink)
	}

	resp, err := json.Marshal(&result)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, "application/json", resp)

}

func (handler *Handler) APIDeleteLinks(ctx *gin.Context) {
	var userID string

	val, ok := ctx.Get("userID")

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok = val.(string)

	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var shorts []string

	if err := json.NewDecoder(ctx.Request.Body).Decode(&shorts); err != nil && !errors.Is(err, io.EOF) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	var batch []store.StoreRecord

	for _, val := range shorts {
		batch = append(batch, store.StoreRecord{
			UUID:      userID,
			ShortLink: val,
		})
	}

	handler.DeleteChan <- batch

	ctx.Status(http.StatusAccepted)
}

func CreateHandler(sStore store.Store, responseAddr string) *Handler {
	handler := &Handler{
		store: sStore,
		options: &options{
			responseAddr: responseAddr,
		},
		DeleteChan: make(chan []store.StoreRecord),
	}

	return handler
}
