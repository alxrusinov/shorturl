package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserLinks - route reteurn all users urls
// /api/user/urls
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
