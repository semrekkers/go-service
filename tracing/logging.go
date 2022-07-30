package tracing

import (
	"net/http"
	"time"

	"github.com/semrekkers/go-service/apih"

	"go.uber.org/zap"
)

type responseCapturer struct {
	http.ResponseWriter
	statusCode int
}

func (c *responseCapturer) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := apih.LoggerFromContext(ctx)
		trace := TraceFromContext(ctx)
		capturer := responseCapturer{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(&capturer, r)
		delta := time.Since(start)

		log.Info("Handled request",
			zap.Int("status", capturer.statusCode),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("trace_id", trace.ID),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", delta),
		)
	})
}
