package service

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/internal/repository"
)

type PaymentService interface {
	CreatePayment(payment *model.Payment) (int64, error)
	GetPaymentById(paymentID int64) (*model.Payment, error)
	GetPaymentByPerson(person string) ([]model.Payment, error)
	UpdatePayment(payment *model.Payment) (int64, error)
	DeletePayment(paymentID int64) (int64, error)
}

type PaymentServiceImpl struct {
	repo repository.PaymentRepository
}

func NewPaymentServiceImpl(repo repository.PaymentRepository) *PaymentServiceImpl {
	return &PaymentServiceImpl{repo: repo}
}

func (p *PaymentServiceImpl) UpdatePayment(payment *model.Payment) (int64, error) {
	return p.repo.Update(payment)
}

func (p *PaymentServiceImpl) DeletePayment(paymentID int64) (int64, error) {
	return p.repo.Delete(paymentID)
}

func (p *PaymentServiceImpl) CreatePayment(payment *model.Payment) (int64, error) {
	return p.repo.Create(payment)
}

func (p *PaymentServiceImpl) GetPaymentById(paymentID int64) (*model.Payment, error) {
	return p.repo.GetById(paymentID)
}

func (p *PaymentServiceImpl) GetPaymentByPerson(person string) ([]model.Payment, error) {
	return p.repo.GetByPerson(person)
}
