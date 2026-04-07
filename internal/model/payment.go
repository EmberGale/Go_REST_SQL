package model

import "time"

type Payment struct {
	Id       int       `json:"id" db:"id"`
	Person   string    `json:"person" db:"person"`
	Amount   int       `json:"amount" db:"amount"`
	Currency string    `json:"currency" db:"currency"`
	Time     time.Time `json:"time" db:"time"`
}
