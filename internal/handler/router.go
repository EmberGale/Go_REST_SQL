package handler

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NewRouter создаёт новый HTTP роутер с middleware логирования
func NewRouter(handler *PaymentHandler, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	// Middleware для логирования запросов
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}

	// Регистрируем маршруты
	mux.HandleFunc("POST /payment/", handler.Create)
	mux.HandleFunc("PUT /payment/", handler.Update)
	mux.HandleFunc("DELETE /payment/", handler.Delete)
	mux.HandleFunc("GET /payment/byId", handler.GetById)
	mux.HandleFunc("GET /payment/byPerson", handler.GetByPerson)

	return loggingMiddleware(mux)
}

// responseWriter обёртка для перехвата статуса ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
