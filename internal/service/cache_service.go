package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
)

func (o *PaymentServiceImpl) getFromCache(ctx context.Context, key string) (interface{}, error) {
	data, err := o.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PaymentServiceImpl) WarmCache(ctx context.Context) error {
	payments, err := p.repo.GetAllPayments(ctx)
	if err != nil {
		return err
	}
	for _, payment := range payments {
		key := fmt.Sprintf("payment_id:%d", payment.ID)
		p.cache.Set(ctx, key, payment, 1*time.Minute)
	}
	return nil
}

func (p *PaymentServiceImpl) GetWithSingleFlight(ctx context.Context, paymentID int64) (*model.Payment, error) {
	singleFlight.Do(ctx, fmt.Sprintf("payment_id:%d", paymentID), func() (interface{}, error) {
		return p.GetPaymentById(paymentID)
	})
	return payment, err
}