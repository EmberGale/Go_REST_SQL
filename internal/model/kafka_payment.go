package model

import "time"

type KafkaPaymentEvent struct {
	ID          int64         `json:"id"`
	Amount      float64       `json:"amount"`
	Currency    string        `json:"currency"`
	Status      PaymentStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	EventType   string        `json:"event_type"`
	AggregateID int64         `json:"aggregate_id"`
}
