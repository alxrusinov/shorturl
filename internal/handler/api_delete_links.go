package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/gin-gonic/gin"
)

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

	var batch []model.StoreRecord

	for _, val := range shorts {
		batch = append(batch, model.StoreRecord{
			UUID:      userID,
			ShortLink: val,
		})
	}

	handler.DeleteChan <- batch

	ctx.Status(http.StatusAccepted)
}
