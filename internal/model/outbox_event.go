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
	ID          string       `json:"id"`
	PaymentID   string       `json:"payment_id"`
	Status      OutboxStatus `json:"status"`
	Payload     string       `json:"payload"` // JSON string of the event
	Attempts    int          `json:"attempts"`
	NextRetryAt time.Time    `json:"next_retry_at"`
	CreatedAt   time.Time    `json:"created_at"`
}
