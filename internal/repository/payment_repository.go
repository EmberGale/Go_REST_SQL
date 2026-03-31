package repository

import "GoRestSQL/internal/model"

// PaymentRepository определяет интерфейс для работы с платежами
type PaymentRepository interface {
	Create(payment *model.Payment) (int64, error)
	GetById(id int64) (*model.Payment, error)
	GetByPerson(person string) ([]model.Payment, error)
	Update(payment *model.Payment) (int64, error)
	Delete(id int64) (int64, error)
}
