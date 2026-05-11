package model

import "time"

type PaymentStatus string

const (
	PAYMENT_STATUS_PENDING PaymentStatus = "PENDING"
	PAYMENT_STATUS_CREATED PaymentStatus = "CREATED"
	PAYMENT_STATUS_FAILED  PaymentStatus = "FAILED"
)

type Payment struct {
	Id       int64         `json:"id" db:"id"`
	Person   string        `json:"person" db:"person"`
	Amount   float64       `json:"amount" db:"amount"`
	Currency string        `json:"currency" db:"currency"`
	Status   PaymentStatus `json:"status" db:"status"`
	Time     time.Time     `json:"time" db:"time"`
	UserID   int64         `json:"user_id" db:"user_id"`
}
