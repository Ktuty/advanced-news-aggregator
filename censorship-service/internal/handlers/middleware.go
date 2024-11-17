package handlers

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

// Middleware для добавления сквозного идентификатора запроса
func (api *API) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		r.Header.Set("X-Request-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), "X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware для логирования запросов
func (api *API) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(&rec, r)
		duration := time.Since(start).Milliseconds()

		requestID := r.Context().Value("X-Request-ID").(string)
		ipAddress := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		statusCode := rec.statusCode

		log.Printf("Request ID: %s, IP: %s, Method: %s, Path: %s, Status: %d, Duration: %dms",
			requestID, ipAddress, method, path, statusCode, duration)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
