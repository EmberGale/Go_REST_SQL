package repository

import "GoRestSQL/internal/model"

type PaymentRepository interface {
	Create(payment *model.Payment) (int64, error)
	GetByPerson(person string) ([]model.Payment, error)
	Update(payment *model.Payment) (int64, error)
	Delete(id int) (int64, error)
}
