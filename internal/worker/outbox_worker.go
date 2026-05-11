package worker

import (
	"GoRestSQL/internal/repository"
	"context"
	"sync"

	"go.uber.org/zap"
)
import "GoRestSQL/pkg/kafka"

type OutboxWorker struct {
	repo        repository.PaymentRepository
	producer    kafka.Producer
	workerCount int
	log         *zap.Logger
}

func NewOutboxWorker(repo repository.PaymentRepository, prod kafka.Producer, workerCount int, log *zap.Logger) *OutboxWorker {
	return &OutboxWorker{
		repo:        repo,
		producer:    prod,
		workerCount: workerCount,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	w.log.Info("Starting outbox workers", zap.Int("Workers number", w.workerCount))
	var wg sync.WaitGroup

	for i := 0; i < w.workerCount; i++ {
		wg.Add(1)
		go w.runWorker(ctx, i, &wg)
	}

	wg.Wait()
	w.log.Info("Finished outbox workers", zap.Int("Workers number", w.workerCount))
}

func (w *OutboxWorker) runWorker(ctx context.Context, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()
	w.log.Info("Starting worker", zap.Int("WorkerID", workerID))

	for {
		select {
		case <-ctx.Done():
			w.log.Info("Stopping worker", zap.Int("WorkerID", workerID))
			return
		default:
			events, err := w.repo

		}
	}
}
