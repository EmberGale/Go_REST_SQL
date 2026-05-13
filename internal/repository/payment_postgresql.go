package repository

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/db"
	"context"
	"encoding/json"
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

	// Payment
	var id int64
	query := `INSERT INTO payments (person, amount, currency, status, user_id, time) 
		VALUES (:person, :amount, :currency, :status, :user_id, :time) RETURNING id`
	payment.Time = time.Now()
	rows, err := tx.NamedQuery(query, payment)

	if err != nil {
		return 0, err
	}
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	// Outbox
	query_outbox := `INSERT INTO outbox_events (payment_id, event_type, payload, status, attempts, 
                           next_retry_at, created_at) 
	VALUES (:payment_id, :event_type, :payload, :status, :attempts, :next_retry_at, 
	        :created_at) RETURNING id`

	payloadJSON, err := json.Marshal(payment)

	outbox_event := model.OutboxEvent{
		ID:          0, // Ignore
		PaymentID:   payment.Id,
		Status:      model.OUTBOX_STATUS_PENDING,
		Payload:     string(payloadJSON),
		Attempts:    0,
		NextRetryAt: time.Time{},
		CreatedAt:   time.Time{},
	}

	rows, err = tx.NamedQuery(query_outbox, outbox_event)
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}
