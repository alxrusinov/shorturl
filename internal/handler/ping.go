package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) Ping(ctx *gin.Context) {
	err := handler.store.Ping()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)

}
