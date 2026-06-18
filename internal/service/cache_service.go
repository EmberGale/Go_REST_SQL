package service

import (
	"context"
	"fmt"
)

func (p *PaymentServiceImpl) WarmCache(ctx context.Context) error {
	p
	key := fmt.Sprintf("payment_id:%d", paymentID)
}
