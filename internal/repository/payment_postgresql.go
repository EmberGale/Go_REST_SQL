package repository

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/db"
	"time"
)

// PostgreSQLPaymentRepository реализует интерфейс PaymentRepository
type PostgreSQLPaymentRepository struct {
	DB db.DB
}

// NewPostgreSQLPaymentRepository создаёт новый репозиторий с подключением к БД
func NewPostgreSQLPaymentRepository(database db.DB) *PostgreSQLPaymentRepository {
	return &PostgreSQLPaymentRepository{DB: database}
}

func (p *PostgreSQLPaymentRepository) Create(payment *model.Payment) (int64, error) {
	var id int64
	query := `INSERT INTO payments (person, amount, currency, time) VALUES ($1, $2, $3, $4) RETURNING id`
	err := p.DB.QueryRow(query, payment.Person, payment.Amount, payment.Currency, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PostgreSQLPaymentRepository) GetById(id int64) (*model.Payment, error) {
	var payment model.Payment
	err := p.DB.QueryRow("SELECT id, person, amount, currency, time FROM payments WHERE id = $1", id).
		Scan(&payment.Id, &payment.Person, &payment.Amount, &payment.Currency, &payment.Time)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *PostgreSQLPaymentRepository) GetByPerson(person string) ([]model.Payment, error) {
	rows, err := p.DB.Query("SELECT id, person, amount, currency, time FROM payments WHERE person = $1", person)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.Id, &payment.Person, &payment.Amount, &payment.Currency, &payment.Time); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (p *PostgreSQLPaymentRepository) Update(payment *model.Payment) (int64, error) {
	query := `UPDATE payments SET person = $1, amount = $2, currency = $3, time = $4 WHERE id = $5`
	res, err := p.DB.Exec(query, payment.Person, payment.Amount, payment.Currency, time.Now(), payment.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (p *PostgreSQLPaymentRepository) Delete(id int64) (int64, error) {
	res, err := p.DB.Exec("DELETE FROM payments WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
