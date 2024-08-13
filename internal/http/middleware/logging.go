package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status        int
	headerWritten bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.headerWritten {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.headerWritten = true
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info("REQ", "method", r.Method, "path", r.URL.EscapedPath())
			ww := wrapResponseWriter(w)
			next.ServeHTTP(ww, r)
			logger.Info("RES", "code", ww.Status(), "duration", time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}
