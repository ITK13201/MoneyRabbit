package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// bodyWriter wraps gin.ResponseWriter to capture the response body for logging.
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// For JSON requests, read and restore the body so the handler can still read it.
		var requestBody string
		ct := c.ContentType()
		if strings.HasPrefix(ct, "application/json") {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			requestBody = string(bodyBytes)
		}

		slog.InfoContext(c.Request.Context(), "http.request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("request_id", c.GetString("request_id")),
			slog.Group("extra",
				"query", c.Request.URL.RawQuery,
				"content_type", ct,
				"body", requestBody,
			),
		)

		// Wrap ResponseWriter to capture response body.
		bw := &bodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = bw

		c.Next()

		slog.InfoContext(c.Request.Context(), "http.response",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("request_id", c.GetString("request_id")),
			slog.Group("extra",
				"status", c.Writer.Status(),
				"latency_ms", time.Since(start).Milliseconds(),
				"body", bw.body.String(),
			),
		)
	}
}
