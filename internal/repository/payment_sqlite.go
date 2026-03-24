package repository

import (
	"GoRestSQL/internal/model"
	"database/sql"
	"time"
)

type SqlitePaymentRepository struct {
	db *sql.DB
}

func NewSqlitePaymentRepository(db *sql.DB) PaymentRepository {
	return &SqlitePaymentRepository{db: db}
}

func (r *SqlitePaymentRepository) Create(payment *model.Payment) (int64, error) {
	const query = `INSERT INTO payments (person, amount, currency, time) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, payment.Person, payment.Amount, payment.Currency, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *SqlitePaymentRepository) GetByPerson(person string) ([]model.Payment, error) {
	const query = `SELECT * FROM payments WHERE person = ?`

	rows, err := r.db.Query(query, person)
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

	result, err := r.db.Exec(query, payment.Person, payment.Amount, payment.Currency, time.Now(), payment.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *SqlitePaymentRepository) Delete(id int) (int64, error) {
	const query = `DELETE FROM payments WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
