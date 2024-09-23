package gateway

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"
)

func (g *Gateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		writer := &statusCodeWriter{ResponseWriter: w}
		next.ServeHTTP(writer, r)

		g.logger.Printf("%s %s : %d %s",
			r.Method, r.RequestURI,
			writer.statusCode,
			time.Since(startTime).String())
	})
}

type statusCodeWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Hijack is implemented to allow clients reach server from different origin.
// NB: Just for example allowing local clients to reach server, should take precaution in real environment.
func (w *statusCodeWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (g *Gateway) logError(err error) {
	g.logger.Printf("Error: %v\n", err)
}
