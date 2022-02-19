package middleware

import (
	"io"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func SetLogger(isTerm bool) gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
			logger := zerolog.New(out)
			if isTerm {
				logger = logger.Output(
					zerolog.ConsoleWriter{
						Out:     out,
						NoColor: false,
					},
				)
			}
			logger = logger.With().
				Timestamp().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("ip", c.ClientIP()).
				Dur("latency", latency).
				Str("user_agent", c.Request.UserAgent()).
				Logger()
			return logger
		}),
	)
}
