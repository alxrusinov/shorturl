package handler

import (
	"net/http"

	"github.com/alxrusinov/shorturl/internal/netutils"
	"github.com/gin-gonic/gin"
)

func (handler Handler) Stats(ctx *gin.Context) {
	ip := ctx.Request.Header.Get("X-Real-IP")

	trusted, err := netutils.CheckSubnet(handler.options.trustedSubnet, ip)

	if !trusted || err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	res, err := handler.store.GetStat()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)

}
