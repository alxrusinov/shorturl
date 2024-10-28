package handler

import (
	"net/http"

	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetOriginalLink(ctx *gin.Context) {
	id := ctx.Param("id")
	defer ctx.Request.Body.Close()

	links := &model.StoreRecord{
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
