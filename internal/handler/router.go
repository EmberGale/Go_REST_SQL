package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// NewRouter создаёт новый HTTP роутер с middleware логирования
func NewRouter(handler *PaymentHandler, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	// Используем стандартные middleware chi
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Middleware для логирования запросов
	r.Use(func(next http.Handler) http.Handler {
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
	})

	// Регистрируем маршруты с chi
	r.Route("/payment", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Get("/byId", handler.GetById)
		r.Get("/byPerson", handler.GetByPerson)
		r.Get("/payment/{id}/inCurrency")
	})

	return r
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
