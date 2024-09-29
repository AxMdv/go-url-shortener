package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/AxMdv/go-url-shortener/pkg/logger"
)

type (
	// responseData is structure for storing information about the response
	responseData struct {
		status int
		size   int
	}
	// loggingResponseWriter implements http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write writes the response using the original http.ResponseWriter.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader writes the status code using the original http.ResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging logs information about request and response
// and returns a new http.Handler
func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h(&lw, r)
		duration := time.Since(start)
		logger.Log.Info("Incoming HTTP request:",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("duration", duration.String()),
		)
		logger.Log.Info("Response to HTTP request:",

			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}
	return http.HandlerFunc(logFn)
}
