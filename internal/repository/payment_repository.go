package repository

import (
	"GoRestSQL/internal/model"
	"context"
)

// PaymentRepository определяет интерфейс для работы с платежами
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) (int64, error)
	GetById(ctx context.Context, id int64) (*model.Payment, error)
	GetByPerson(ctx context.Context, person string) ([]model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	CreateWithOutbox(ctx context.Context, payment *model.Payment) (int64, error)
}
