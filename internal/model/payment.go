package model

import "time"

type Payment struct {
	Id       int       `json:"id"`
	Person   string    `json:"person"`
	Amount   int       `json:"amount"`
	Currency string    `json:"currency"`
	Time     time.Time `json:"time"`
}
