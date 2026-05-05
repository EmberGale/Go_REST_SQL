package kafka

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/config"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type KafkaProducer interface {
	SendPaymentCreated(message *PaymentCreatedMessage) error
	Close() error
}

type PaymentCreatedMessage struct {
	EventID   string        `json:"event_id"`
	EventType string        `json:"event_type"`
	Timestamp string        `json:"timestamp"`
	Payment   model.Payment `json:"payment"`
}

type Producer struct {
	producer *kafka.Producer
	topic    string
	logger   *zap.Logger
}

func NewKafkaProducer(cfg *config.KafkaConfig, topic string, logger *zap.Logger) (*Producer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
	}

	p, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, fmt.Errorf("init kafka producer: %w", err)
	}

	return &Producer{producer: p, topic: topic}, nil
}

func (p *Producer) SendPaymentCreated(message *PaymentCreatedMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal payment created message to json: %w", err)
	}

	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)
}

func (p *Producer) Close() {
	p.producer.Close()
}
