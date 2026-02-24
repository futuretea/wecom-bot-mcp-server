package http

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/futuretea/wecom-bot-mcp-server/pkg/core/logging"
)

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		logging.Debug("%s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.headerWritten {
		return
	}
	lrw.statusCode = code
	lrw.headerWritten = true
	lrw.ResponseWriter.WriteHeader(code)
}

// Write writes the response body data. If WriteHeader has not been called yet,
// it implicitly sets the status to 200 OK, matching net/http stdlib behavior.
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.WriteHeader(http.StatusOK)
	return lrw.ResponseWriter.Write(b)
}

func (lrw *loggingResponseWriter) Flush() {
	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := lrw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
