package repository

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/db"
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// PostgreSQLPaymentRepository реализует интерфейс PaymentRepository
type PostgreSQLPaymentRepository struct {
	DB db.DB
}

// NewPostgreSQLPaymentRepository создаёт новый репозиторий с подключением к БД
func NewPostgreSQLPaymentRepository(database db.DB) *PostgreSQLPaymentRepository {
	return &PostgreSQLPaymentRepository{DB: database}
}

func (p *PostgreSQLPaymentRepository) Create(ctx context.Context, payment *model.Payment) (int64, error) {
	var id int64
	query := `INSERT INTO payments (person, amount, currency, time) VALUES (:person, :amount, :currency, :time) RETURNING id`
	payment.Time = time.Now()
	rows, err := p.DB.NamedQuery(query, payment)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (p *PostgreSQLPaymentRepository) GetById(ctx context.Context, id int64) (*model.Payment, error) {
	var payment model.Payment
	query := "SELECT id, person, amount, currency, time FROM payments WHERE id = $1"
	if err := p.DB.Get(&payment, query, id); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *PostgreSQLPaymentRepository) GetByPerson(ctx context.Context, person string) ([]model.Payment, error) {
	query := "SELECT id, person, amount, currency, time FROM payments WHERE person = $1"
	var payments []model.Payment
	if err := p.DB.Select(&payments, query, person); err != nil {
		return nil, err
	}
	return payments, nil
}

func (p *PostgreSQLPaymentRepository) Update(ctx context.Context, payment *model.Payment) (int64, error) {
	query := `UPDATE payments SET person = :person, amount = :amount, currency = :currency, time = :time WHERE id = :id`
	payment.Time = time.Now()
	result, err := p.DB.NamedExec(query, payment)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (p *PostgreSQLPaymentRepository) Delete(ctx context.Context, id int64) (int64, error) {
	result, err := p.DB.Exec("DELETE FROM payments WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (p *PostgreSQLPaymentRepository) CreateWithOutbox(ctx context.Context, payment *model.Payment) (int64, error) {
	tx, err := p.DB.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	payment.Time = time.Now()

	// 1. Insert Payment
	queryPayment := `INSERT INTO payments (person, amount, currency, status, user_id, time) 
		VALUES (:person, :amount, :currency, :status, :user_id, :time) RETURNING id`

	q1, args1, err := sqlx.Named(queryPayment, payment)
	if err != nil {
		return 0, err
	}

	var id int64
	// tx.Rebind converts '?' placeholders to '$1, $2' for PostgreSQL
	if err := tx.QueryRowx(tx.Rebind(q1), args1...).Scan(&id); err != nil {
		return 0, err
	}

	// 2. Insert Outbox
	payloadJSON, _ := json.Marshal(payment)

	outboxEvent := model.OutboxEvent{
		PaymentID: id, // Use the newly generated ID
		Status:    model.OUTBOX_STATUS_PENDING,
		Payload:   string(payloadJSON),
		Attempts:  0,
		//NextRetryAt: time.Time{}, empty value, as it will be set when the first retry is scheduled
		//CreatedAt:   time.Time{}, will be set by the database default value (CURRENT_TIMESTAMP)
		Topic: "payment_created",
	}

	queryOutbox := `INSERT INTO outbox_events (payment_id, payload, status, attempts, topic) 
		VALUES (:payment_id, :payload, :status, :attempts, :topic) RETURNING id`

	q2, args2, err := sqlx.Named(queryOutbox, outboxEvent)
	if err != nil {
		return 0, err
	}

	var outboxID int64
	if err := tx.QueryRowx(tx.Rebind(q2), args2...).Scan(&outboxID); err != nil {
		return 0, err
	}

	return outboxID, tx.Commit()
}
