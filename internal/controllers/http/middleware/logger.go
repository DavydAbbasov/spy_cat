package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		responseWriter := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseWriter

		c.Next()

		logResponse(c, responseWriter, startTime)
	}
}

func logResponse(c *gin.Context, w *bodyLogWriter, startTime time.Time) {

	logger := log.With().
		Str("type", "requestResponse").
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status_code", c.Writer.Status()).
		Int64("duration_ms", time.Since(startTime).Milliseconds()).
		Int("response_size", c.Writer.Size()).
		Logger()

	var requestBody string
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	if requestBody != "" {
		if isJSON(requestBody) {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(requestBody), "", "  "); err == nil {
				logger = logger.With().Str("body", prettyJSON.String()).Logger()
			} else {
				logger = logger.With().Str("body", requestBody).Logger()
			}
		} else {
			logger = logger.With().Str("body", requestBody).Logger()
		}

		responseBody := w.body.String()

		if responseBody != "" {
			if isJSON(responseBody) {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, []byte(responseBody), "", "  "); err == nil {
					logger = logger.With().Str("body", prettyJSON.String()).Logger()
				} else {
					logger = logger.With().Str("body", responseBody).Logger()
				}
			} else {
				logger = logger.With().Str("body", responseBody).Logger()
			}
		}

		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Error().Msg("Outgoing response")
		case status >= 400:
			logger.Warn().Msg("Outgoing response")
		default:
			logger.Info().Msg("Outgoing response")
		}
	}
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
