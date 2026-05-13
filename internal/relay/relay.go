package relay

import (
	"GoRestSQL/pkg/kafka"
	"context"
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
}

func NewRelay(config RelayConfig, kafkaProduce *kafka.Producer, log *zap.Logger) *Relay {
	ctx, cancel := context.WithCancel(context.Background())
	return &Relay{
		config:        config,
		logger:        zap.NewNop(),
		ctx:           ctx,
		cancel:        cancel,
		kafkaProducer: kafkaProduce,
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
		}
		case<- ticker.C:
			r.processBatch(id)
	}
}

func (r *Relay) processBatch(id int) {
	defer r.wg.Done()
	r.logger.Info("process batch" + strconv.Itoa(id))

}
