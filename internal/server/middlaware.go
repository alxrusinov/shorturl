package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var zippingContentType = [2]string{"text/html", "application/json"}

var zipFormat = "gzip"

func checkContentType(values []string) bool {
	var result bool
	for _, value := range values {
		if value == zippingContentType[0] || value == zippingContentType[1] {
			result = true
			break
		}
	}
	return result
}

func checkGzip(values []string) bool {
	var result bool
	for _, value := range values {
		if value == zipFormat {
			result = true
			break
		}
	}
	return result
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

func loggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
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

func compressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := c.Request.Header.Values("Content-Encoding")
		acceptEncoding := c.Request.Header.Values("Accept-Encoding")

		contentType := c.Request.Header.Values("Content-Type")

		if !checkContentType(contentType) {
			c.Next()
			return
		}

		if checkGzip(contentEncoding) {
			rawContent, err := gzip.NewReader(c.Request.Body)

			if err != nil && err != io.EOF {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			defer rawContent.Close()

			c.Request.Body = rawContent
			c.Request.Header.Set("Content-Encoding", "identity")

		}

		if checkGzip(acceptEncoding) {
			gz := gzip.NewWriter(c.Writer)
			c.Header("Content-Encoding", "gzip")
			c.Writer = &gzipWriter{c.Writer, gz}
		}

		c.Next()

	}
}
