package relay

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/db"
	"GoRestSQL/pkg/kafka"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type RelayConfig struct {
	WorkerCount   int           `mapstructure:"worker_count"`
	BatchSize     int           `mapstructure:"batch_size"`
	PollInterval  time.Duration `mapstructure:"poll_interval"`
	MaxAttempts   int           `mapstructure:"max_attempts"`
	BaseDelay     time.Duration `mapstructure:"base_delay"`
	MaxDelay      time.Duration `mapstructure:"max_delay"`
	JitterPercent float64       `mapstructure:"jitter_percent"`
}
type Relay struct {
	config        RelayConfig
	kafkaProducer *kafka.Producer
	logger        *zap.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	db            db.DB
}

func NewRelay(config RelayConfig, kafkaProduce *kafka.Producer, log *zap.Logger, db db.DB) *Relay {
	ctx, cancel := context.WithCancel(context.Background())
	return &Relay{
		config:        config,
		logger:        log,
		ctx:           ctx,
		cancel:        cancel,
		kafkaProducer: kafkaProduce,
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
	defer r.wg.Done()
	r.logger.Info("process batch" + strconv.Itoa(id))

	query_tasks := `
	UPDATE outbox_events
	SET status = 'processing'
	WHERE id IN (
		SELECT id 
		FROM outbox_events 
		WHERE status = 'pending' AND next_retry_at <= NOW()
		ORDER BY created_at ASC
		LIMIT 10
		FOR UPDATE SKIP LOCKED
	)
	RETURNING id, payment_id, payload, attempts, topic;
	`
	var events []model.OutboxEvent
	err := r.db.Select(&events, query_tasks)
	if err != nil {
		r.logger.Error("failed to fetch events", zap.Error(err))
	}

	for _, event := range events {
		msg := kafka.Message{
			Topic: event.Topic,
			Key:   []byte(fmt.Sprintf("%d", event.ID)),
			Value: []byte(event.Payload),
		}

		err := r.kafkaProducer.SendMessage(r.ctx, msg)
		if err != nil {
			r.logger.Error("failed to send message", zap.Error(err))
		}
	}

}
