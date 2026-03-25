package repository

import (
	"GoRestSQL/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgreSQLPaymentRepository struct {
	db *sqlx.DB
}

func NewPostgreSQLPaymentRepository() (PaymentRepository, *sqlx.DB, error) {
	dsn := "user=postgres password=123 dbname=payments sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		return nil, nil, err
	}

	if _, err = db.Exec(schemaSQL); err != nil {
		return nil, nil, err
	}
	return &SqlitePaymentRepository{db: db}, db, nil
}

func (p *PostgreSQLPaymentRepository) Create(payment *model.Payment) (int64, error) {
	res, err := p.db.Exec("INSERT INTO payments (person, amount, currency, time) VALUES ?, ?, ?, ?)", payment.Person, payment.Amount, payment.Currency, time.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (p *PostgreSQLPaymentRepository) GetById(id int64) (*model.Payment, error) {
	var payment model.Payment
	err := p.db.Select(&payment, "SELECT * FROM payments WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *PostgreSQLPaymentRepository) GetByPerson(person string) ([]model.Payment, error) {
	var payments []model.Payment
	err := p.db.Select(&payments, "SELECT * FROM payments WHERE person = $1", person)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (p *PostgreSQLPaymentRepository) Update(payment *model.Payment) (int64, error) {
	const query = `UPDATE payments SET person = ?, amount = ?, currency = ? WHERE id = ?`
	res, err := p.db.Exec(query, payment.Person, payment.Amount, payment.Currency, time.Now(), payment.Id)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (p *PostgreSQLPaymentRepository) Delete(id int64) (int64, error) {
	res, err := p.db.Exec("DELETE FROM payments WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
