package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alxrusinov/shorturl/internal/model"
)

// APIShortenBatch - route adds urls by batch
// /api/shorten/batch
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

	var content []*model.StoreRecord

	if err := json.NewDecoder(ctx.Request.Body).Decode(&content); err != nil && !errors.Is(err, io.EOF) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	for _, val := range content {
		shortenURL, err := handler.Generator.GenerateRandomString()

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
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusCreated, "application/json", resp)

}
