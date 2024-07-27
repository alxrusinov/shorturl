package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Middleware func() gin.HandlerFunc

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
