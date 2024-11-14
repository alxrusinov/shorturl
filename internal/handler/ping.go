package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping - handler for pinging store
func (handler *Handler) Ping(ctx *gin.Context) {
	err := handler.store.Ping()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)

}
