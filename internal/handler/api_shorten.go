package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/customerrors"
	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/gin-gonic/gin"
)

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

	links := &model.StoreRecord{
		ShortLink:    shortenURL,
		OriginalLink: content.URL,
		UUID:         userID,
	}

	res, err := handler.store.SetLink(links)

	dbErr := &customerrors.DuplicateValueError{}

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
