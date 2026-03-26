package repository

import (
	"GoRestSQL/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type SqlitePaymentRepository struct {
	DB *sqlx.DB
}

func NewSqlitePaymentRepository() (*SqlitePaymentRepository, error) {
	db, err := sqlx.Connect("sqlite", "payments.db")
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err = db.Exec(schemaSQL); err != nil {
		return nil, err
	}
	return &SqlitePaymentRepository{DB: db}, nil
}

func (r *SqlitePaymentRepository) GetById(id int64) (*model.Payment, error) {
	const query = `SELECT * FROM payments WHERE id = ?`

	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	var payment model.Payment
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&payment.Person, &payment.Amount, &payment.Currency, &payment.Time); err != nil {
			return nil, err
		}
		return &payment, nil
	}
	return nil, nil
}

func (r *SqlitePaymentRepository) Create(payment *model.Payment) (int64, error) {
	const query = `INSERT INTO payments (person, amount, currency, time) VALUES (?, ?, ?, ?)`

	result, err := r.DB.Exec(query, payment.Person, payment.Amount, payment.Currency, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *SqlitePaymentRepository) GetByPerson(person string) ([]model.Payment, error) {
	const query = `SELECT * FROM payments WHERE person = ?`

	rows, err := r.DB.Query(query, person)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.Person, &payment.Amount, &payment.Currency, &payment.Time); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *SqlitePaymentRepository) Update(payment *model.Payment) (int64, error) {
	const query = `UPDATE payments SET person = ?, amount = ?, currency = ? WHERE id = ?`

	result, err := r.DB.Exec(query, payment.Person, payment.Amount, payment.Currency, time.Now(), payment.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *SqlitePaymentRepository) Delete(id int64) (int64, error) {
	const query = `DELETE FROM payments WHERE id = ?`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
