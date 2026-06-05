package service

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/internal/repository"
	"GoRestSQL/pkg/config"
	"GoRestSQL/pkg/kafka"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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
	ctx           context.Context
	repo          repository.PaymentRepository
	kafkaProducer kafka.Producer
	log           *zap.SugaredLogger
	cache         *redis.Client
	cfg           *config.Config
}

func NewPaymentServiceImpl(ctx context.Context, repo repository.PaymentRepository,
	kafkaProducer kafka.Producer, log *zap.SugaredLogger, cache *redis.Client) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		log:           log,
		ctx:           ctx,
		cache:         cache,
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
	p.log.Info("Get by id: %d", paymentID)

	// Try cache
	key := fmt.Sprintf("payment_id:%d", paymentID)
	if data, err := p.cache.Get(p.ctx, key).Result(); err == nil {
		var payment model.Payment
		_ = json.Unmarshal([]byte(data), &payment)
		p.log.Info("Get by id:%d -> from cache", paymentID)
		return &payment, nil
	}
	// Cache miss
	payment, err := p.repo.GetById(p.ctx, paymentID)

	// Store in cache
	bytes, _ := json.Marshal(payment)
	p.log.Info("Get by id:%d -> save to cache", paymentID)
	cacheErr := p.cache.Set(p.ctx, key, bytes, 1*time.Minute).Err()
	if cacheErr != nil {
		p.log.Info("Get by id:%d -> failed to cache: %s", paymentID, cacheErr)
	}
	return payment, err
}

func (p *PaymentServiceImpl) GetPaymentByPerson(person string) ([]model.Payment, error) {
	return p.repo.GetByPerson(p.ctx, person)
}
