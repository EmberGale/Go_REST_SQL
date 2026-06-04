package relay

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/config"
	"GoRestSQL/pkg/db"
	"GoRestSQL/pkg/kafka"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Relay struct {
	config        config.RelayConfig
	kafkaProducer kafka.Producer
	logger        *zap.SugaredLogger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	db            db.DB
}

func NewRelay(config config.RelayConfig, kafkaProducer kafka.Producer, log *zap.SugaredLogger, db db.DB) *Relay {
	ctx, cancel := context.WithCancel(context.Background())
	return &Relay{
		config:        config,
		logger:        log,
		ctx:           ctx,
		cancel:        cancel,
		kafkaProducer: kafkaProducer,
		db:            db,
	}
}

// Start worker pool
func (r *Relay) Start() {
	r.logger.Info("start relay")

	for i := 0; i < r.config.WorkerCount; i++ {
		r.wg.Add(1)
		go r.worker(i)
	}
}

func (r *Relay) Stop() {
	r.logger.Info("stop relay")
	r.cancel()
	r.wg.Wait()
}

func (r *Relay) worker(id int) {
	defer r.wg.Done()
	r.logger.Info("worker started" + strconv.Itoa(id))

	ticker := time.NewTicker(r.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("worker stopping" + strconv.Itoa(id))
			return

		case <-ticker.C:
			r.processBatch(id)
		}
	}
}

func (r *Relay) processBatch(id int) {
	r.logger.Info("process batch" + strconv.Itoa(id))

	queryTasks := `
	UPDATE outbox_events
	SET status = $1
	WHERE id IN (
		SELECT id 
		FROM outbox_events 
		WHERE status = $2 AND next_retry_at <= NOW()
		ORDER BY created_at ASC
		LIMIT 10
		FOR UPDATE SKIP LOCKED
	)
	RETURNING id, payment_id, payload, attempts, topic;
	`

	var events []model.OutboxEvent
	err := r.db.Select(&events, queryTasks,
		model.OUTBOX_STATUS_PROCESSING,
		model.OUTBOX_STATUS_PENDING)
	if err != nil {
		r.logger.Error("failed to fetch events", zap.Error(err))
		return
	}

	for _, event := range events {
		msg := kafka.Message{
			Topic: event.Topic,
			Key:   []byte(fmt.Sprintf("%d", event.ID)),
			Value: []byte(event.Payload),
		}

		// Status to change to after kafka attempt
		var newStatus model.OutboxStatus
		err := r.kafkaProducer.SendMessage(r.ctx, msg)
		if err != nil {
			r.logger.Error("failed to send message", zap.Error(err))

			if event.Attempts >= r.config.MaxAttempts {
				newStatus = model.OUTBOX_STATUS_FAILED
			} else {
				newStatus = model.OUTBOX_STATUS_PENDING
			}
		} else {
			newStatus = model.OUTBOX_STATUS_SUCCESS
		}

		_, err = r.db.Exec(
			"UPDATE outbox_events SET status = $1, attempts = $2 WHERE id = $3",
			newStatus,
			event.Attempts,
			event.ID,
		)

		if err != nil {
			r.logger.Error("failed to update event", zap.Int64("id", event.ID), zap.Error(err))
		}
	}

}
