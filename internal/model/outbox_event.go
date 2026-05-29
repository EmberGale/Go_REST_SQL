package model

import "time"

type OutboxStatus string

const (
	OUTBOX_STATUS_PENDING    OutboxStatus = "PENDING"
	OUTBOX_STATUS_PROCESSING OutboxStatus = "PROCESSING"
	OUTBOX_STATUS_SUCCESS    OutboxStatus = "SUCCESS"
	OUTBOX_STATUS_FAILED     OutboxStatus = "FAILED"
)

type OutboxEvent struct {
	ID          int64        `json:"id" db:"id"`
	PaymentID   int64        `json:"payment_id" db:"payment_id"`
	EventType   string       `json:"event_type" db:"event_type"`
	Status      OutboxStatus `json:"status" db:"status"`
	Payload     string       `json:"payload" db:"payload"`
	Attempts    int          `json:"attempts" db:"attempts"`
	NextRetryAt time.Time    `json:"next_retry_at" db:"next_retry_at"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	Topic       string       `json:"topic" db:"topic"`
}
