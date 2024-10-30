package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alxrusinov/shorturl/internal/customerrors"
	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/model"
)

// GetShortLink - route adds url
// /
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

	links := &model.StoreRecord{
		ShortLink:    shortenURL,
		OriginalLink: originURL,
		UUID:         userID,
	}

	res, err := handler.store.SetLink(links)

	dbErr := &customerrors.DuplicateValueError{}

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
