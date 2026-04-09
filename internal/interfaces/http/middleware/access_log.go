package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logAny, ok := c.Get(ContextKeyLogger)
		if !ok {
			return
		}
		log, ok := logAny.(*slog.Logger)
		if !ok {
			return
		}
		log.Info("requisição concluída",
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.Int("bytes", c.Writer.Size()),
		)
	}
}
