package kafka

import (
	"context"
)

// Message represents the data to be sent
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Partition int32
}

// Producer is the abstract interface
type Producer interface {
	SendMessage(ctx context.Context, msg Message) error
	Close() error
}
