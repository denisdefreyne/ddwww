package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type httpResponseWriter struct {
	http.ResponseWriter
	HTTPStatus   int
	ResponseSize int
}

func (w *httpResponseWriter) WriteHeader(status int) {
	w.HTTPStatus = status

	w.ResponseWriter.WriteHeader(status)
}

func (w *httpResponseWriter) Write(b []byte) (int, error) {
	w.ResponseSize = len(b)

	return w.ResponseWriter.Write(b)
}

func noCachingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w2 := httpResponseWriter{w, 0, 0}

		next.ServeHTTP(&w2, r)

		// Extract IP from RemoteAddr (which also contains the port)
		remoteIp := strings.Split(r.RemoteAddr, ":")[0]

		// Get status (use 200 by default)
		httpStatus := w2.HTTPStatus
		if httpStatus == 0 {
			httpStatus = 200
		}

		log.Printf("%s %s %s %s \"%s %s %s\" %d %d\n",
			remoteIp,
			"-", // identity
			"-", // user
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			httpStatus,
			w2.ResponseSize,
		)
	})
}
