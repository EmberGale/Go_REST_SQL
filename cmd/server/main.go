package main

import (
	"GoRestSQL/internal/handler"
	"GoRestSQL/internal/repository"
	"GoRestSQL/internal/service"
	"GoRestSQL/pkg/config"
	"GoRestSQL/pkg/db"
	"GoRestSQL/pkg/kafka"
	"GoRestSQL/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	// Создаём логгер
	log, err := logger.New(cfg.Logger.Level)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}
	defer log.Sync()

	log.Info("application starting", zap.String("port", cfg.Server.Port))

	// Подключаемся к БД и применяем миграции
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	log.Info("database connected and migrations applied")

	// Kafka
	kafkaProducer, err := kafka.NewKafkaProducer(&cfg.Kafka, "Payment", log)
	if err != nil {
		log.Fatal("failed to create kafka producer", zap.Error(err))
	}

	// Создаём слои приложения
	paymentRepo := repository.NewPostgreSQLPaymentRepository(database)
	paymentService := service.NewPaymentServiceImpl(paymentRepo, kafkaProducer)
	paymentHandler := handler.NewPaymentHandler(paymentService, log)
	router := handler.NewRouter(paymentHandler, log)

	// Создаём HTTP сервер
	addr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Info("server starting", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	// Ждём сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server stopped")
}
