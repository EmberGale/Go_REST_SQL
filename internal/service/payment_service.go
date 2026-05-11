package service

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/internal/repository"
	"GoRestSQL/pkg/kafka"

	"go.uber.org/zap"
)

type PaymentService interface {
	CreatePayment(payment *model.Payment) (int64, error)
	GetPaymentById(paymentID int64) (*model.Payment, error)
	GetPaymentByPerson(person string) ([]model.Payment, error)
	UpdatePayment(payment *model.Payment) (int64, error)
	DeletePayment(paymentID int64) (int64, error)
}

type PaymentServiceImpl struct {
	repo          repository.PaymentRepository
	kafkaProducer kafka.Producer
	log           *zap.Logger
}

func NewPaymentServiceImpl(repo repository.PaymentRepository, kafkaProducer kafka.Producer, log *zap.Logger) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		log:           log,
	}
}

func (p *PaymentServiceImpl) UpdatePayment(payment *model.Payment) (int64, error) {
	return p.repo.Update(payment)
}

func (p *PaymentServiceImpl) DeletePayment(paymentID int64) (int64, error) {
	return p.repo.Delete(paymentID)
}

func (p *PaymentServiceImpl) CreatePayment(payment *model.Payment) (int64, error) {
	resp, err := p.repo.Create(payment)
	if err != nil {
		return 0, err
	}

	return resp, err
}

func (p *PaymentServiceImpl) GetPaymentById(paymentID int64) (*model.Payment, error) {
	return p.repo.GetById(paymentID)
}

func (p *PaymentServiceImpl) GetPaymentByPerson(person string) ([]model.Payment, error) {
	return p.repo.GetByPerson(person)
}
