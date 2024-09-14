package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"time"

	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Middlewares struct{}

const UserCookie = "user_cookie"

func checkContentType(values []string) bool {
	var zippingContentType = map[string]struct{}{"text/html": {}, "application/json": {}}

	for _, value := range values {
		if _, ok := zippingContentType[value]; ok {
			return true
		}
	}
	return false
}

func checkGzip(values []string) bool {
	const zipFormat = "gzip"

	for _, value := range values {
		if value == zipFormat {
			return true
		}
	}
	return false
}

type Middleware func() gin.HandlerFunc

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func (middlwares *Middlewares) LoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		uri := c.Request.RequestURI
		method := c.Request.Method

		c.Next()

		size := c.Writer.Size()
		status := c.Writer.Status()

		duration := time.Since(start)

		logger.Info().Str("uri", uri).Str("method", method).Int("status", status).Dur("duration", duration).Int("size", size)

	}
}

func (middlwares *Middlewares) CompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := c.Request.Header.Values("Content-Encoding")
		acceptEncoding := c.Request.Header.Values("Accept-Encoding")

		contentType := c.Request.Header.Values("Content-Type")

		if !checkContentType(contentType) && c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		if checkGzip(contentEncoding) {
			rawContent, err := gzip.NewReader(c.Request.Body)

			if err != nil && err != io.EOF {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			defer rawContent.Close()

			c.Request.Body = rawContent
			c.Request.Header.Set("Content-Encoding", "identity")

		}

		c.Next()

		if checkGzip(acceptEncoding) {
			gz := gzip.NewWriter(c.Writer)
			c.Writer.Header().Set("Content-Encoding", "gzip")
			c.Writer = &gzipWriter{c.Writer, gz}
		}

	}
}

func (middlwares *Middlewares) CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()
		method := c.Request.Method

		if method == http.MethodPost {
			userID, err := c.Cookie(UserCookie)

			if err != nil {
				userID, err = generator.GenerateUserID()

				if err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				c.SetCookie(UserCookie, userID, 60*60*24, "/", "localhost", false, true)

			}

			c.Set("userID", userID)

			c.Next()

			return
		}

		if method == http.MethodGet {
			if fullPath == "/api/user/urls" {
				userID, err := c.Cookie(UserCookie)

				if err != nil {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}

				c.Set("userID", userID)
			}

			c.Next()

			return
		}

		if method == http.MethodDelete {
			userID, err := c.Cookie(UserCookie)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Set("userID", userID)

			c.Next()

			return
		}

	}
}
