package producer

import (
	"context"
	"fmt"

	"GoRestSQL/pkg/kafka"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type saramaProducer struct {
	syncProducer sarama.SyncProducer
	log          *zap.SugaredLogger
}

// NewSaramaProducer creates a concrete Kafka producer backed by Sarama.
func NewSaramaProducer(brokers []string, config *sarama.Config, logger *zap.SugaredLogger) (kafka.Producer, error) {
	if config == nil {
		config = sarama.NewConfig()
		config.Producer.Return.Successes = true
	}

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Return the error so the caller can handle it gracefully
		return nil, fmt.Errorf("failed to create Sarama producer: %w", err)
	}

	return &saramaProducer{
		syncProducer: p,
		log:          logger,
	}, nil
}

func (s *saramaProducer) SendMessage(ctx context.Context, msg kafka.Message) error {
	// Check if context was canceled before blocking on Sarama
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sMsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}

	_, _, err := s.syncProducer.SendMessage(sMsg)
	if err != nil {
		s.log.Error("failed to send Kafka message", zap.Error(err))
	}
	return err
}

func (s *saramaProducer) Close() error {
	return s.syncProducer.Close()
}
