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

		logRespWriter := &logResponseWriter{ResponseWriter: w}
		next.ServeHTTP(logRespWriter, r)

		g.logger.Printf("%s %s : %d %s",
			r.Method, r.RequestURI,
			logRespWriter.statusCode,
			time.Since(startTime).String())
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (g *Gateway) logError(err error) {
	g.logger.Printf("Error: %v\n", err)
}
