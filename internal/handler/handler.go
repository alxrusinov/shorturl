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
	store   store.Store
	options *options
}

const UserCookie = "user_cookie"

func (handler *Handler) GetShortLink(ctx *gin.Context) {
	userID, err := ctx.Cookie(UserCookie)

	if err != nil {
		userID, err = generator.GenerateUserID()

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.SetCookie(UserCookie, userID, 60*60*24, "/", "localhost", false, true)
	}

	body, _ := io.ReadAll(ctx.Request.Body)
	originURL := string(body)

	shortenURL, err := generator.GenerateRandomString(10)

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
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Header("Location", res.OriginalLink)
	ctx.Status(http.StatusTemporaryRedirect)
}

func (handler *Handler) APIShorten(ctx *gin.Context) {
	userID, err := ctx.Cookie(UserCookie)

	if err != nil {
		userID, err = generator.GenerateUserID()

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.SetCookie(UserCookie, userID, 60*60*24, "/", "localhost", false, true)
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

	shortenURL, err = generator.GenerateRandomString(10)

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
	userID, err := ctx.Cookie(UserCookie)

	if err != nil {
		userID, err = generator.GenerateUserID()
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.SetCookie(UserCookie, userID, 60*60*24, "/", "localhost", false, true)
	}

	var content []*store.StoreRecord

	if err := json.NewDecoder(ctx.Request.Body).Decode(&content); err != nil && !errors.Is(err, io.EOF) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer ctx.Request.Body.Close()

	for _, val := range content {
		shortenURL, err := generator.GenerateRandomString(10)

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
	userID, err := ctx.Cookie(UserCookie)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
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

func (handler *Handler) APIDeleteLinks(ctx *gin.Context) {}

func CreateHandler(store store.Store, responseAddr string) *Handler {
	handler := &Handler{
		store: store,
		options: &options{
			responseAddr: responseAddr,
		},
	}

	return handler
}
