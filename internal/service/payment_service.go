package service

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/internal/repository"
	"GoRestSQL/pkg/kafka"
	"context"

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
	log           *zap.SugaredLogger
	ctx           context.Context
}

func NewPaymentServiceImpl(ctx context.Context, repo repository.PaymentRepository, kafkaProducer kafka.Producer, log *zap.SugaredLogger) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		log:           log,
		ctx:           ctx,
	}
}

func (p *PaymentServiceImpl) UpdatePayment(payment *model.Payment) (int64, error) {
	return p.repo.Update(p.ctx, payment)
}

func (p *PaymentServiceImpl) DeletePayment(paymentID int64) (int64, error) {
	return p.repo.Delete(p.ctx, paymentID)
}

func (p *PaymentServiceImpl) CreatePayment(payment *model.Payment) (int64, error) {
	return p.repo.CreateWithOutbox(p.ctx, payment)
}

func (p *PaymentServiceImpl) GetPaymentById(paymentID int64) (*model.Payment, error) {
	return p.repo.GetById(p.ctx, paymentID)
}

func (p *PaymentServiceImpl) GetPaymentByPerson(person string) ([]model.Payment, error) {
	return p.repo.GetByPerson(p.ctx, person)
}
