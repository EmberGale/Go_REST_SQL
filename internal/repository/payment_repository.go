package repository

import "GoRestSQL/internal/model"

type PaymentRepository interface {
	Create(payment *model.Payment) (int64, error)
	GetById(id int64) (*model.Payment, error)
	GetByPerson(person string) ([]model.Payment, error)
	Update(payment *model.Payment) (int64, error)
	Delete(id int64) (int64, error)
}

const (
	schemaSQL = `CREATE TABLE IF NOT EXISTS payments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	person TEXT,
	amount REAL,
	currency TEXT,
	time INTEGER,
	);`

	dbfile = "payments.db"
)
